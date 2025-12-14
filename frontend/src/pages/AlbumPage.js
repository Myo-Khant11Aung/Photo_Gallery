import { useParams, useNavigate } from "react-router-dom";
import PhotoCard from "../components/PhotoCard";
import { useState, useEffect, useCallback } from "react";
import "lightgallery/css/lightgallery.css";
import "lightgallery/css/lg-zoom.css";
import "lightgallery/css/lg-thumbnail.css";
import LightGallery from "lightgallery/react";
import lgZoom from "lightgallery/plugins/zoom";
import lgThumbnail from "lightgallery/plugins/thumbnail";

const API = process.env.REACT_APP_API;

function AlbumPage() {
  const { date } = useParams();
  const [images, setImages] = useState([]);
  const token = localStorage.getItem("token");
  const refreshImages = useCallback(() => {
    fetch(`${API}/api/images`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((data) => setImages(data.filter((img) => img.album_date === date)));
  }, [date, token]);

  function handleDelete(id) {
    setImages((prev) => prev.filter((img) => img.id !== id));
  }

  useEffect(() => {
    refreshImages();
  }, [refreshImages]);

  const navigate = useNavigate();
  return (
    <div className="photo-page">
      <div className="album-toolbar">
        <button className="back-btn" onClick={() => navigate(`/`)}>
          ← Back
        </button>
      </div>

      <div className="album-header">
        <h2 className="album-title">Album for {date}</h2>
      </div>
      <LightGallery
        dynamic={true}
        dynamicEl={images.map((img) => ({
          src: img.url,
          thumb: img.url,
        }))}
        plugins={[lgZoom, lgThumbnail]}
        speed={300}
        onInit={(detail) => {
          window.lightGalleryInstance = detail.instance;
        }}
      ></LightGallery>

      <div className="photo-grid">
        {images.map((image, index) => (
          <PhotoCard
            key={image.id}
            image={image}
            index={index}
            onMemoUpdated={refreshImages}
            onDelete={handleDelete}
          />
        ))}
      </div>
      {/* <button
        type="button"
        className="upload-fab"
        onClick={openPicker}
        title={uploading ? "Uploading…" : "Upload photos"}
        disabled={uploading}
      >
        <svg viewBox="0 0 24 24">
          <path
            d="M12 5v14m-7-7h14"
            stroke="white"
            strokeWidth="2"
            fill="none"
            strokeLinecap="round"
          />
        </svg>
      </button> */}
    </div>
  );
}
export default AlbumPage;

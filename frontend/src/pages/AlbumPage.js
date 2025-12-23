import { useParams, useNavigate } from "react-router-dom";
import PhotoCard from "../components/PhotoCard";
import { useState, useEffect, useCallback } from "react";
import "lightgallery/css/lightgallery.css";
import "lightgallery/css/lg-zoom.css";
import "lightgallery/css/lg-thumbnail.css";
import LightGallery from "lightgallery/react";
import lgZoom from "lightgallery/plugins/zoom";
import lgThumbnail from "lightgallery/plugins/thumbnail";
import UploadHandler from "../components/uploadHandler";

const API = process.env.REACT_APP_API;

function AlbumPage() {
  const { albumId } = useParams();
  const [images, setImages] = useState([]);
  const token = localStorage.getItem("token");
  const refreshImages = useCallback(() => {
    fetch(`${API}/api/albums/${albumId}/images`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((data) => setImages(data || []));
  }, [albumId, token]);

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
        <UploadHandler albumId={albumId} onUploadSuccess={refreshImages} />
      </div>

      <div className="album-header">
        <h2 className="album-title">Album for {albumId}</h2>
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
    </div>
  );
}
export default AlbumPage;

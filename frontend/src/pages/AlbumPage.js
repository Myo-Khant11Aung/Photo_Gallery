import { useParams, useNavigate } from "react-router-dom";
import PhotoCard from "../components/PhotoCard";
import { useState, useEffect, useCallback } from "react";
import "lightgallery/css/lightgallery.css";
import "lightgallery/css/lg-zoom.css";
import "lightgallery/css/lg-thumbnail.css";
import LightGallery from "lightgallery/react";
import lgZoom from "lightgallery/plugins/zoom";
import lgThumbnail from "lightgallery/plugins/thumbnail";

function AlbumPage() {
  const { date } = useParams();
  const [images, setImages] = useState([]);
  const token = localStorage.getItem("token");
  const refreshImages = useCallback(() => {
    fetch("http://localhost:8080/api/images", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((data) => setImages(data.filter((img) => img.album_date === date)));
  }, [date, token]);

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
          />
        ))}
      </div>
    </div>
  );
}
export default AlbumPage;

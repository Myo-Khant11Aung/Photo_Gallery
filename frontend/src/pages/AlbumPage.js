import { useParams, useNavigate } from "react-router-dom";
import PhotoCard from "../components/PhotoCard";
import { useState, useEffect, useCallback } from "react";

function AlbumPage() {
  const { date } = useParams();
  const [images, setImages] = useState([]);

  const refreshImages = useCallback(() => {
    fetch("http://localhost:8080/api/images")
      .then((res) => res.json())
      .then((data) => setImages(data.filter((img) => img.album_date === date)));
  }, [date]);

  useEffect(() => {
    refreshImages();
  }, [refreshImages]);

  const navigate = useNavigate();
  return (
    <div>
      <button onClick={() => navigate(`/`)}>Back</button>
      <h2>Album for {date}</h2>
      <div>
        {images &&
          images.map((image) => (
            <PhotoCard
              key={image.id}
              image={image}
              onMemoUpdated={refreshImages}
            />
          ))}
      </div>
    </div>
  );
}
export default AlbumPage;

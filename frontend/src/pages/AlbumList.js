import { useState, useEffect, useCallback, useRef } from "react";
import { useNavigate } from "react-router-dom";
import LogoutButton from "../components/LogoutButton";
import "../styles.css";
import heic2any from "heic2any";
import UploadHandler from "../components/uploadHandler";

const API = process.env.REACT_APP_API;


function App() {
  const [images, setImages] = useState([]);
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const refreshImages = useCallback(() => {
    fetch(`${API}/api/images`, {
      method: "GET",
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((data) => setImages(data || []));
  }, [token]);
  useEffect(() => {
    refreshImages();
  }, [refreshImages]);
  const album = {};
  images.forEach((image) => {
    const date = image.album_date;
    if (!album[date]) {
      album[date] = [];
    }
    album[date].push(image);
  });

function handleCreateAlbum() {
    // simplest version: "new album" = today's date
    const today = new Date().toISOString().slice(0, 10); // YYYY-MM-DD
    navigate(`/album/${today}`);
  }

  return (
    <div className="album-page">
      <div className="album-toolbar">
        {/* <UploadHandler onUploadSuccess={refreshImages} /> */}
        <LogoutButton />
      </div>

      <div className="album-header">
        <h1 className="album-title">Photo Gallery</h1>
      </div>
      <div className="album-list">
        {Object.entries(album).map(([date, img]) => (
          <div
            key={date}
            className="album-card"
            onClick={() => navigate(`/album/${date}`)}
          >
            <h3 className="album-card-title">{date}</h3>

            <div className="album-media">
              <img src={img[0].url} alt="" className="album-thumbnail" />
            </div>

            <p className="album-meta">{img.length} photo(s)</p>
          </div>
        ))}
      </div>

      {/* Floating "Create Album" circle button */}
      <button
        type="button"
        className="upload-fab"
        onClick={handleCreateAlbum}
        title="Create new album"
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
      </button>
    </div>
  );
}

export default App;

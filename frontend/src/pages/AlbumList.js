import { useState, useEffect, useCallback, useRef } from "react";
import { useNavigate } from "react-router-dom";
import LogoutButton from "../components/LogoutButton";
import "../styles.css";
import heic2any from "heic2any";

const API = process.env.REACT_APP_API;

async function convertIfHeic(file) {
  const lowerName = file.name.toLowerCase();

  // If it's not a HEIC/HEIF file, just return it as-is
  if (
    !lowerName.endsWith(".heic") &&
    !lowerName.endsWith(".heif") &&
    file.type !== "image/heic" &&
    file.type !== "image/heif"
  ) {
    return file;
  }

  // Convert HEIC → JPEG using heic2any
  const convertedBlob = await heic2any({
    blob: file,
    toType: "image/jpeg",
    quality: 0.9,
  });

  // Wrap the Blob back into a File so FormData works normally
  const newName = lowerName.endsWith(".heic")
    ? file.name.replace(/\.heic$/i, ".jpg")
    : lowerName.endsWith(".heif")
    ? file.name.replace(/\.heif$/i, ".jpg")
    : file.name + ".jpg";

  return new File([convertedBlob], newName, { type: "image/jpeg" });
}

function UploadHandler({ onUploadSuccess }) {
  const fileInputRef = useRef(null);
  const [uploading, setUploading] = useState(false);
  const token = localStorage.getItem("token");

  function openPicker() {
    fileInputRef.current?.click();
  }

  async function uploadManyAtOnce(files) {
    const formData = new FormData();

    // Convert each file if needed
    for (const f of files) {
      const processed = await convertIfHeic(f);
      formData.append("image", processed); // backend still reads "image"
    }

    const res = await fetch(`${API}/api/upload`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: formData,
    });

    if (!res.ok) {
      throw new Error(`Upload failed: ${res.status}`);
    }

    return res.json();
  }

  function handleChange(e) {
    const files = Array.from(e.target.files || []);
    if (!files.length || uploading) return;

    setUploading(true);
    uploadManyAtOnce(files)
      .then(() => onUploadSuccess && onUploadSuccess())
      .catch((err) => {
        console.error("Upload Error", err);
        alert("Upload failed. Please try again.");
      })
      .finally(() => {
        setUploading(false);
        // allow re-selecting the same files later
        e.target.value = "";
      });
  }

  return (
    <>
      {/* Hidden multi-file input */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        multiple
        style={{ display: "none" }}
        onChange={handleChange}
      />

      {/* Plus icon button */}
      <button
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
      </button>
    </>
  );
}
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
  return (
    <div className="album-page">
      <div className="album-toolbar">
        <UploadHandler onUploadSuccess={refreshImages} />
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
    </div>
  );
}

export default App;

import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import LogoutButton from "/Users/myokhantaung/Desktop/Photo_Gallery/frontend/src/components/LogoutButton.js";
function UploadHandler({ onUploadSuccess }) {
  const [selectedFile, setSelectedFile] = useState([]);
  function handleFileChange(e) {
    setSelectedFile(e.target.files[0]);
  }
  const token = localStorage.getItem("token");
  function handleUpload(e) {
    e.preventDefault();

    const formData = new FormData();
    formData.append("image", selectedFile);

    fetch("http://localhost:8080/api/upload", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: formData,
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("upload complete", data);
        onUploadSuccess();
      })
      .catch((err) => console.error("Upload Error", err));
  }
  return (
    <div>
      <form onSubmit={handleUpload}>
        <input type="file" onChange={handleFileChange} />
        <button type="submit">Upload</button>
        <LogoutButton />
      </form>
    </div>
  );
}

function App() {
  const [images, setImages] = useState([]);
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const refreshImages = useCallback(() => {
    fetch("http://localhost:8080/api/images", {
      method: "GET",
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((data) => setImages(data));
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
    <div>
      <UploadHandler onUploadSuccess={refreshImages} />
      <div>
        <h1>Photo Gallery</h1>
        <div>
          {Object.entries(album).map(([date, img]) => (
            <div key={date} onClick={() => navigate(`/album/${date}`)}>
              <h3>{date}</h3>
              <img
                src={`http://localhost:8080/images/${encodeURIComponent(
                  img[0].filename
                )}`}
                alt=""
                width="200"
              />
              <p>{img.length} photo(s)</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;

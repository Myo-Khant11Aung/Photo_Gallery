import { useState } from "react";
import "lightgallery/css/lightgallery.css";
import "lightgallery/css/lg-zoom.css";
import "lightgallery/css/lg-thumbnail.css";
import LightGallery from "lightgallery/react";
import lgZoom from "lightgallery/plugins/zoom";
import lgThumbnail from "lightgallery/plugins/thumbnail";

function PhotoCard({ image, onMemoUpdated }) {
  const [inputMemo, setInputMemo] = useState(image.memo || "");
  const [isEditing, setIsEditing] = useState(false);
  const token = localStorage.getItem("token");
  function handleEditingStatus() {
    setIsEditing(true);
  }

  function handleClear() {
    setInputMemo("");
  }

  function handleMemoChange(e) {
    setInputMemo(e.target.value);
  }

  function handleMemoUpload() {
    fetch(`http://localhost:8080/api/photo/${image.id}/memo`, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ memo: inputMemo }),
    })
      .then((res) => res.json())
      .then((data) => {
        setIsEditing(false);
        if (onMemoUpdated) onMemoUpdated();
      })
      .catch((err) => console.error("Failed to update memo", err));
  }
  const src = `http://localhost:8080/images/${encodeURIComponent(
    image.filename
  )}`;

  return (
    <div className="photo-card">
      <a href={src} className="lightbox-item" data-lg-size="1600-1200">
        <img src={src} alt="" />
      </a>
      <div className="photo-card-body">
        {isEditing ? (
          <>
            <textarea
              value={inputMemo}
              onChange={handleMemoChange}
              placeholder="Memo about this photo..."
            />
            <button onClick={handleMemoUpload}>Confirm</button>
            <button onClick={handleClear}>Clear</button>
          </>
        ) : (
          <>
            <p>{image.memo || "No memo yet"}</p>
            <button onClick={handleEditingStatus}>Memo</button>
          </>
        )}
      </div>
    </div>
  );
}

export default PhotoCard;

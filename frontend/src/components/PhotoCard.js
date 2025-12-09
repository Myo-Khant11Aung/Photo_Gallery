import { useState } from "react";
import "lightgallery/css/lightgallery.css";
import "lightgallery/css/lg-zoom.css";
import "lightgallery/css/lg-thumbnail.css";
import LightGallery from "lightgallery/react";
import lgZoom from "lightgallery/plugins/zoom";
import lgThumbnail from "lightgallery/plugins/thumbnail";

const API = process.env.REACT_APP_API;

function PhotoCard({ image, onMemoUpdated, index }) {
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
    fetch(`${API}/api/photo/${image.id}/memo`, {
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
  const src = image.url;

  return (
    <div className="photo-card">
      <div
        className="photo-card-image"
        onClick={() => {
          if (window.lightGalleryInstance) {
            window.lightGalleryInstance.openGallery(index);
          }
        }}
      >
        <img src={src} alt="" />
      </div>
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

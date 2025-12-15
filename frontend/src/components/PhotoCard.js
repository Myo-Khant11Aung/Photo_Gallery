import { useState } from "react";
import "lightgallery/css/lightgallery.css";
import "lightgallery/css/lg-zoom.css";
import "lightgallery/css/lg-thumbnail.css";

const API = process.env.REACT_APP_API;

function PhotoCard({ image, onMemoUpdated, index, onDelete }) {
  const [inputMemo, setInputMemo] = useState(image.memo || "");
  const [isEditing, setIsEditing] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);

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
  function handleDeleteClick(e) {
    e.stopPropagation();

    fetch(`${API}/api/photo/delete/${image.id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Delete failed");
        return res.json();
      })
      .then(() => {
        setMenuOpen(false);
        onDelete?.(image.id);
      })
      .catch((err) => console.error("Failed to delete", err));
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
      <button
        className="photo-menu-btn"
        onClick={(e) => {
          e.stopPropagation();
          setMenuOpen(!menuOpen);
        }}
      >
        ⋮
      </button>

      {menuOpen && (
        <div className="photo-menu" onClick={(e) => e.stopPropagation()}>
          <button className="delete-btn" onClick={handleDeleteClick}>
            Delete Photo
          </button>
        </div>
      )}

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

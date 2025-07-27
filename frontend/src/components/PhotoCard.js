import { useState } from "react";

function PhotoCard({ image, onMemoUpdated }) {
  const [inputMemo, setInputMemo] = useState(image.memo || "");
  const [isEditing, setIsEditing] = useState(false);

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

  return (
    <div className="photo-card">
      <img
        src={`http://localhost:8080/images/${encodeURIComponent(
          image.filename
        )}`}
        alt=""
        width="200"
      />
      {isEditing ? (
        <>
          <br />
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
  );
}

export default PhotoCard;

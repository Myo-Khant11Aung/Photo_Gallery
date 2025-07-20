import { useState, useEffect } from "react";

function PhotoCard({ image }) {
  const [inputMemo, setInputMemo] = useState(image.memo || "");

  function handleMemoChange(e) {
    setInputMemo(e.target.value);
  }

  function handleMemoUpload() {
    console.log(`Confirm memo upload ${image.id}: ${inputMemo}`);
  }

  return (
    <div>
      <img
        src={`http://localhost:8080/images/${encodeURIComponent(
          image.filename
        )}`}
        alt=""
        width="200"
      />
      <br />
      <input
        type="text"
        value={inputMemo}
        onChange={handleMemoChange}
        placeholder="Memo about this photo..."
        size={30}
      />
      <button onClick={handleMemoUpload}>Confirm</button>
      <p>{image.memo}</p>
    </div>
  );
}

function App() {
  const [images, setImages] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/images")
      .then((res) => res.json())
      .then((data) => setImages(data));
  }, []);

  const [selectedFile, setSelectedFile] = useState([]);
  function handleFileChange(e) {
    setSelectedFile(e.target.files[0]);
  }

  function handleUpload(e) {
    e.preventDefault();

    const formData = new FormData();
    formData.append("image", selectedFile);

    fetch("http://localhost:8080/api/upload", {
      method: "POST",
      body: formData,
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("upload complete", data);
        fetch("http://localhost:8080/api/images")
          .then((res) => res.json())
          .then((data) => setImages(data));
      })
      .catch((err) => console.error("Upload Error", err));
  }
  return (
    <div>
      <form onSubmit={handleUpload}>
        <input type="file" onChange={handleFileChange} />
        <button type="submit">Upload</button>
      </form>
      {images &&
        images.map((image) => <PhotoCard key={image.filename} image={image} />)}
    </div>
  );
}

export default App;

import { useState, useEffect } from "react";

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
        images.map((image) => (
          <div key={image.filename}>
            <h1>{image.title}</h1>
            <img
              src={`http://localhost:8080/images/${encodeURIComponent(
                image.filename
              )}`}
              alt={image.title}
              width="200"
            />
          </div>
        ))}
    </div>
  );
}

export default App;

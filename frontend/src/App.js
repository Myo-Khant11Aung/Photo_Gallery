import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AlbumList from "./pages/AlbumList";
import AlbumPage from "./pages/AlbumPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<AlbumList />}></Route>
        <Route path="/album/:date" element={<AlbumPage />}></Route>
      </Routes>
    </Router>
  );
}
export default App;

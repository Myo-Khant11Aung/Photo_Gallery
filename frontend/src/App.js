import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AlbumList from "./pages/AlbumList";
import AlbumPage from "./pages/AlbumPage";
import Login from "./pages/Login";
import Register from "./pages/Register";
import PrivateRoute from "./components/PrivateRoute";

function App() {
  return (
    <Router>
      <Routes>
        <Route element={<PrivateRoute />}>
          <Route path="/" element={<AlbumList />}></Route>
          <Route path="/album/:date" element={<AlbumPage />}></Route>
        </Route>
        <Route path="/login" element={<Login />}></Route>
        <Route path="/register" element={<Register />}></Route>
      </Routes>
    </Router>
  );
}
export default App;

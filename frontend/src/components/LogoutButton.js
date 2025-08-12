import { useNavigate } from "react-router-dom";
import "../styles.css";
function LogoutButton() {
  const navigate = useNavigate();

  function handleLogout() {
    localStorage.removeItem("token");
    navigate("/login");
  }

  return (
    <button onClick={handleLogout} className="back-btn">
      Log Out
    </button>
  );
}

export default LogoutButton;

import { Navigate, Outlet } from "react-router-dom";
import { useState, useEffect } from "react";

const API = process.env.REACT_APP_API;

function PrivateRoute() {
  const [checked, setChecked] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const token = localStorage.getItem("token");
  // return token ? <Outlet /> : <Navigate to="/login" replace />;
  useEffect(() => {
    if (!token) {
      setChecked(true);
      setIsAuthenticated(false);
      return;
    } else {
      fetch(`${API}/api/verifyToken`, {
        method: "GET",
        headers: { Authorization: `Bearer ${token}` },
      })
        .then((res) => {
          if (res.status === 200) {
            setIsAuthenticated(true);
          } else {
            localStorage.removeItem("token");
            setIsAuthenticated(false);
          }
        })
        .catch((err) => {
          console.error("Token verification failed:", err);
          setIsAuthenticated(false);
        })
        .finally(() => {
          setChecked(true);
        });
    }
  }, [token]);

  if (!checked) {
    return <div>Loading...</div>;
  }

  // 2. Done checking
  if (isAuthenticated) {
    return <Outlet />;
  } else {
    return <Navigate to="/login" replace />;
  }
}

export default PrivateRoute;

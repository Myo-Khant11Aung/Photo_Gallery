import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles.css";

const API = process.env.REACT_APP_API;

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const token = localStorage.getItem("token");
  const navigate = useNavigate();
  function handleSubmit(e) {
    e.preventDefault();

    fetch(`${API}/api/login`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password_hash: password }),
    })
      .then((res) => {
        if (!res.ok) throw new Error("Login failed");
        return res.json();
      })
      .then((data) => {
        console.log("Login success:", data);
        localStorage.setItem("token", data.token);
        navigate("/");
      })
      .catch((err) => {
        if (err.message.includes("403")) {
          alert(
            "You do not have access yet. Please wait until you're assigned a wall."
          );
        }
      });
  }

  return (
    <div className="auth-page">
      <div className="auth-container">
        <form onSubmit={handleSubmit} className="auth-form">
          <h2>Login</h2>
          <input
            type="email"
            placeholder="Email"
            className="auth-input"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            className="auth-input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <div className="button-group">
            <button type="submit" className="auth-button">
              Log In
            </button>
            <button
              type="button"
              onClick={() => navigate("/register")}
              className="auth-button"
            >
              Sign Up
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Login;

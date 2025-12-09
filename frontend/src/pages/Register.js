import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import "../styles.css";

const API = process.env.REACT_APP_API;

function Register() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const token = localStorage.getItem("token");
  const navigate = useNavigate();
  function handleSubmit(e) {
    e.preventDefault();

    fetch(`${API}/api/register`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username,
        email,
        password_hash: password,
      }),
    })
      .then((res) => {
        if (!res.ok) throw new Error("Registration failed");
        return res.json();
      })
      .then((data) => {
        console.log("Registration successful:", data);
        navigate("/login");
      })
      .catch((err) => {
        console.error("Error:", err.message);
      });
  }

  return (
    <div className="auth-page">
      <div className="auth-container">
        <form onSubmit={handleSubmit} className="auth-form">
          <h2>Register</h2>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="auth-input"
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="auth-input"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="auth-input"
          />
          <button type="submit" className="auth-button">
            Register
          </button>
          <div>
            <Link to={"/login"}>Already have an account? Login</Link>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Register;

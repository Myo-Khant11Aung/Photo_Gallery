import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const token = localStorage.getItem("token");
  const navigate = useNavigate();
  function handleSubmit(e) {
    e.preventDefault();

    fetch("http://localhost:8080/api/login", {
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
    <form onSubmit={handleSubmit}>
      <h2>Login</h2>
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <button type="submit" onClick={handleSubmit}>
        Log In
      </button>
      <div>
        <Link to={"/register"}>Register?</Link>
      </div>
    </form>
  );
}

export default Login;

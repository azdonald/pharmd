import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Register() {
  const [form, setForm] = useState({
    firstName: "", lastName: "", email: "", orgName: "", password: "", confirmPassword: "",
  });
  const [error, setError] = useState("");
  const { register, isLoading } = useAuth();
  const navigate = useNavigate();

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (form.password !== form.confirmPassword) {
      setError("Passwords do not match");
      return;
    }
    try {
      await register(form.email, form.password, form.orgName, form.firstName, form.lastName);
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    }
  };

  return (
    <div className="auth-page">
      <form className="auth-form" onSubmit={handleSubmit}>
        <h1>Register</h1>
        {error && <p className="error">{error}</p>}
        <label>First Name<input value={form.firstName} onChange={update("firstName")} required /></label>
        <label>Last Name<input value={form.lastName} onChange={update("lastName")} required /></label>
        <label>Email<input type="email" value={form.email} onChange={update("email")} required /></label>
        <label>Organisation<input value={form.orgName} onChange={update("orgName")} required /></label>
        <label>Password<input type="password" value={form.password} onChange={update("password")} required minLength={8} /></label>
        <label>Confirm Password<input type="password" value={form.confirmPassword} onChange={update("confirmPassword")} required /></label>
        <button type="submit" disabled={isLoading}>{isLoading ? "Loading..." : "Register"}</button>
        <p className="auth-link">Already have an account? <Link to="/login">Login</Link></p>
      </form>
    </div>
  );
}

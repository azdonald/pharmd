import { useState } from "react";
import * as authApi from "../api/auth";

export default function ChangePassword() {
  const [form, setForm] = useState({ oldPassword: "", newPassword: "", confirmPassword: "" });
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setMessage("");
    if (form.newPassword !== form.confirmPassword) {
      setError("Passwords do not match");
      return;
    }
    try {
      await authApi.changePassword(form.oldPassword, form.newPassword);
      setMessage("Password changed successfully");
      setForm({ oldPassword: "", newPassword: "", confirmPassword: "" });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to change password");
    }
  };

  return (
    <div className="form-page">
      <h1>Change Password</h1>
      {message && <p className="success">{message}</p>}
      {error && <p className="error">{error}</p>}
      <form onSubmit={handleSubmit}>
        <label>Current Password<input type="password" value={form.oldPassword} onChange={update("oldPassword")} required /></label>
        <label>New Password<input type="password" value={form.newPassword} onChange={update("newPassword")} required minLength={8} /></label>
        <label>Confirm New Password<input type="password" value={form.confirmPassword} onChange={update("confirmPassword")} required /></label>
        <button type="submit">Change Password</button>
      </form>
    </div>
  );
}

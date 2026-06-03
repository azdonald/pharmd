import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getUser, createUser, updateUser, assignUserRole } from "../api/users";
import { listRoles, type Role } from "../api/roles";

export default function UserForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new";

  const [form, setForm] = useState({ first_name: "", last_name: "", email: "" });
  const [roleId, setRoleId] = useState("");
  const [roles, setRoles] = useState<Role[]>([]);
  const [saving, setSaving] = useState(false);

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value });

  useEffect(() => {
    listRoles().then(res => setRoles(res.data)).catch(console.error);
    if (!isNew && id) {
      getUser(id).then(u => setForm({ first_name: u.first_name, last_name: u.last_name, email: u.email })).catch(console.error);
    }
  }, [id, isNew]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (isNew) {
        const user = await createUser({ ...form, role_id: roleId || undefined });
        navigate(`/users/${user.id}`);
      } else {
        await updateUser(id!, form);
        if (roleId) {
          await assignUserRole(id!, roleId);
        }
        navigate("/users");
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-page">
      <h1>{isNew ? "New User" : "Edit User"}</h1>
      <form onSubmit={handleSubmit}>
        <label>First Name<input value={form.first_name} onChange={update("first_name")} required /></label>
        <label>Last Name<input value={form.last_name} onChange={update("last_name")} required /></label>
        <label>Email<input type="email" value={form.email} onChange={update("email")} required /></label>
        <label>Role<select value={roleId} onChange={e => setRoleId(e.target.value)}>
          <option value="">None</option>
          {roles.map(r => <option key={r.id} value={r.id}>{r.name}</option>)}
        </select></label>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</button>
      </form>
    </div>
  );
}

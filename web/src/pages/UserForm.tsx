import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getUser, createUser, updateUser, assignUserRole } from "../api/users";
import { listRoles, type Role } from "../api/roles";
import { listLocations, type Location } from "../api/locations";
import { useToast } from "../context/ToastContext";

export default function UserForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new";

  const [form, setForm] = useState({ first_name: "", last_name: "", email: "" });
  const [roleId, setRoleId] = useState("");
  const [locationId, setLocationId] = useState("");
  const [superAdmin, setSuperAdmin] = useState(false);
  const [roles, setRoles] = useState<Role[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [saving, setSaving] = useState(false);
  const { showToast } = useToast();

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value });

  useEffect(() => {
    Promise.all([
      listRoles(),
      listLocations(1, 100),
    ]).then(([roleRes, locRes]) => {
      setRoles(roleRes.data);
      setLocations(locRes.data);
    }).catch(console.error);
    if (!isNew && id) {
      getUser(id).then(u => {
        setForm({ first_name: u.first_name, last_name: u.last_name, email: u.email });
        if (u.location_id) {
          setLocationId(u.location_id);
        } else {
          setSuperAdmin(true);
        }
      }).catch(console.error);
    }
  }, [id, isNew]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isNew && !superAdmin && !locationId) { showToast("Select a location", "error"); return; }
    setSaving(true);
    try {
      const locValue = superAdmin ? "" : locationId;
      if (isNew) {
        const user = await createUser({ ...form, role_id: roleId || undefined, location_id: locValue || undefined });
        showToast("User created successfully");
        navigate(`/app/users/${user.id}`);
      } else {
        await updateUser(id!, { ...form, location_id: locValue || undefined });
        if (roleId) {
          await assignUserRole(id!, roleId);
        }
        showToast("User updated successfully");
        navigate("/app/users");
      }
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Save failed", "error");
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
        <label>Location<select value={locationId} onChange={e => setLocationId(e.target.value)} disabled={superAdmin} required={!superAdmin}>
          <option value="">{superAdmin ? "All Locations" : isNew ? "Select location" : "All Locations"}</option>
          {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
        </select></label>
        <label style={{ marginTop: 8 }}>
          <input type="checkbox" checked={superAdmin} onChange={e => { setSuperAdmin(e.target.checked); if (e.target.checked) setLocationId(""); }} />
          {" "}Super admin / All locations
        </label>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</button>
      </form>
    </div>
  );
}

import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getRole, createRole, updateRole, getRolePermissions, setRolePermissions, type Role } from "../api/roles";
import { listPermissions, type Permission } from "../api/permissions";

export default function RoleForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new";

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [selectedPerms, setSelectedPerms] = useState<string[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    listPermissions().then(res => setPermissions(res.data)).catch(console.error);
    if (!isNew && id) {
      getRole(id).then(r => { setName(r.name); setDescription(r.description || ""); }).catch(console.error);
      getRolePermissions(id).then(r => setSelectedPerms(r.permission_ids)).catch(console.error);
    }
  }, [id, isNew]);

  const togglePerm = (permId: string) => {
    setSelectedPerms(prev =>
      prev.includes(permId) ? prev.filter(p => p !== permId) : [...prev, permId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      let role: Role;
      if (isNew) {
        role = await createRole({ name, description: description || undefined, permission_ids: selectedPerms });
      } else {
        role = await updateRole(id!, { name: name || undefined, description: description || undefined });
        await setRolePermissions(id!, selectedPerms);
      }
      navigate(`/roles/${role.id}`);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-page">
      <h1>{isNew ? "New Role" : "Edit Role"}</h1>
      <form onSubmit={handleSubmit}>
        <label>Name<input value={name} onChange={e => setName(e.target.value)} required /></label>
        <label>Description<input value={description} onChange={e => setDescription(e.target.value)} /></label>
        <fieldset>
          <legend>Permissions</legend>
          {permissions.map(p => (
            <label key={p.id} className="checkbox-label">
              <input type="checkbox" checked={selectedPerms.includes(p.id)} onChange={() => togglePerm(p.id)} />
              {p.name} <code>({p.slug})</code>
            </label>
          ))}
        </fieldset>
        <button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</button>
      </form>
    </div>
  );
}

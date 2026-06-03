import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listRoles, deleteRole, type Role } from "../api/roles";
import { useToast } from "../context/ToastContext";

export default function Roles() {
  const { showToast } = useToast();
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);

  const load = () => {
    setLoading(true);
    listRoles().then(res => setRoles(res.data)).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this role?")) return;
    try {
      await deleteRole(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Roles</h1>
        <Link to="/roles/new" className="btn">New Role</Link>
      </div>
      <table>
        <thead>
          <tr><th>Name</th><th>Slug</th><th>Description</th><th>System</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {roles.map(r => (
            <tr key={r.id}>
              <td>{r.name}</td>
              <td><code>{r.slug}</code></td>
              <td>{r.description}</td>
              <td>{r.is_system ? "Yes" : "No"}</td>
              <td>
                <Link to={`/roles/${r.id}`}>Edit</Link>
                {!r.is_system && <button onClick={() => handleDelete(r.id)}>Delete</button>}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

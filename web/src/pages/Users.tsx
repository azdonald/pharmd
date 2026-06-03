import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listUsers, deleteUser, type User } from "../api/users";

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  const load = () => {
    setLoading(true);
    listUsers().then(res => setUsers(res.data)).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Deactivate this user?")) return;
    try {
      await deleteUser(id);
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Users</h1>
        <Link to="/users/new" className="btn">New User</Link>
      </div>
      <table>
        <thead>
          <tr><th>Name</th><th>Email</th><th>Active</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id}>
              <td>{u.first_name} {u.last_name}</td>
              <td>{u.email}</td>
              <td><span className={`badge ${u.is_active ? "badge-active" : "badge-inactive"}`}>{u.is_active ? "Active" : "Inactive"}</span></td>
              <td>
                <Link to={`/users/${u.id}`} className="action-link">Edit</Link>
                <button onClick={() => handleDelete(u.id)} className="action-link action-link-danger">Deactivate</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

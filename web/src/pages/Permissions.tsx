import { useEffect, useState } from "react";
import { listPermissions, type Permission } from "../api/permissions";

export default function Permissions() {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listPermissions()
      .then(res => setPermissions(res.data))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <h1>Permissions</h1>
      <table>
        <thead>
          <tr><th>ID</th><th>Name</th><th>Slug</th><th>Description</th></tr>
        </thead>
        <tbody>
          {permissions.map(p => (
            <tr key={p.id}>
              <td>{p.id}</td>
              <td>{p.name}</td>
              <td><code>{p.slug}</code></td>
              <td>{p.description}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

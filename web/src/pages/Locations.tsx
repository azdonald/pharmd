import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listLocations, deleteLocation, type Location } from "../api/locations";
import { useToast } from "../context/ToastContext";

export default function Locations() {
  const { showToast } = useToast();
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(true);

  const load = () => {
    setLoading(true);
    listLocations().then(res => setLocations(res.data)).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this location?")) return;
    try {
      await deleteLocation(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Locations</h1>
        <Link to="/locations/new" className="btn">New Location</Link>
      </div>
      <table>
        <thead>
          <tr><th>Name</th><th>City</th><th>State</th><th>Phone</th><th>Active</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {locations.map(l => (
            <tr key={l.id}>
              <td>{l.name}</td>
              <td>{l.city}</td>
              <td>{l.state}</td>
              <td>{l.phone}</td>
              <td>{l.is_active ? "Yes" : "No"}</td>
              <td>
                <Link to={`/locations/${l.id}`}>Edit</Link>
                <button onClick={() => handleDelete(l.id)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listLocations, deleteLocation, type Location } from "../api/locations";
import { useToast } from "../context/ToastContext";
import { Icon, PageHeader, Panel } from "../components/AdminComponents";


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

  return (
    <div>
      <PageHeader
        title="Locations"
        description="Manage pharmacy branches and facilities"
        actions={
          <Link to="/app/locations/new" className="btn-sky-action">
            <Icon name="add" className="mr-2" />New Location
          </Link>
        }
      />

      <Panel>
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading locations...</p></div>
          ) : locations.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No locations found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">City</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">State</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Phone</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {locations.map(l => (
                  <tr key={l.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-semibold text-on-surface">{l.name}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{l.city}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{l.state}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{l.phone}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        l.is_active ? "bg-secondary-container/20 text-secondary" : "bg-outline-variant/20 text-on-surface-variant"
                      }`}>{l.is_active ? "Active" : "Inactive"}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-2">
                        <Link to={`/app/locations/${l.id}`} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">Edit</Link>
                        <button onClick={() => handleDelete(l.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </Panel>
    </div>
  );
}

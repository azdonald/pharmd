import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listRoles, deleteRole, type Role } from "../api/roles";
import { useToast } from "../context/ToastContext";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

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

  return (
    <div>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Roles</h2>
          <p className="text-body-lg text-on-surface-variant">Define user roles and access levels</p>
        </div>
        <Link to="/app/roles/new" className="flex items-center px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container shadow-md transition-all">
          <Icon name="add" className="mr-2" />New Role
        </Link>
      </div>

      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading roles...</p></div>
          ) : roles.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No roles found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Slug</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Description</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">System</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {roles.map(r => (
                  <tr key={r.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-semibold text-on-surface">{r.name}</td>
                    <td className="px-6 py-4"><code className="text-sm bg-surface-container-high px-2 py-0.5 rounded">{r.slug}</code></td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{r.description}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        r.is_system ? "bg-tertiary-container/20 text-tertiary" : "bg-outline-variant/20 text-on-surface-variant"
                      }`}>{r.is_system ? "Yes" : "No"}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-2">
                        <Link to={`/app/roles/${r.id}`} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">Edit</Link>
                        {!r.is_system && (
                          <button onClick={() => handleDelete(r.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}

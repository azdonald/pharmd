import { useEffect, useState } from "react";
import { listPermissions, type Permission } from "../api/permissions";
import { PageHeader, Panel } from "../components/AdminComponents";

export default function Permissions() {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listPermissions()
      .then(res => setPermissions(res.data))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <PageHeader
        title="Permissions"
        description="System permissions and access control rules"
      />

      <Panel>
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading permissions...</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">ID</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Slug</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {permissions.map(p => (
                  <tr key={p.id} className="hover:bg-surface-container-high/20 transition-colors">
                    <td className="px-6 py-4 font-data-mono text-data-mono text-on-surface-variant">{p.id}</td>
                    <td className="px-6 py-4 font-semibold text-on-surface">{p.name}</td>
                    <td className="px-6 py-4"><code className="text-sm bg-surface-container-high px-2 py-0.5 rounded">{p.slug}</code></td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.description}</td>
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

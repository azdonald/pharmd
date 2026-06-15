import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listSuppliers, deleteSupplier, type Supplier } from "../api/suppliers";
import { useToast } from "../context/ToastContext";
import { Icon, PageHeader, Panel } from "../components/AdminComponents";


export default function Suppliers() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const { showToast } = useToast();
  const limit = 20;

  const load = () => {
    setLoading(true);
    listSuppliers(page, limit, search).then(res => {
      setSuppliers(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, search]);

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this supplier?")) return;
    try {
      await deleteSupplier(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      <PageHeader
        title="Suppliers"
        description="Manage vendor and manufacturer relationships"
        actions={
          <Link to="/app/suppliers/new" className="btn-sky-action">
            <Icon name="add" className="mr-2" />New Supplier
          </Link>
        }
      />

      <div className="mb-8 max-w-md">
        <div className="relative">
          <Icon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant" />
          <input value={search} onChange={e => { setSearch(e.target.value); setPage(1); }}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest py-2 pl-10 pr-4 text-sm outline-none focus:ring-2 focus:ring-primary"
            placeholder="Search suppliers..." type="text" />
        </div>
      </div>

      <Panel>
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading suppliers...</p></div>
          ) : suppliers.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No suppliers found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Contact</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Phone</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Email</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {suppliers.map(s => (
                  <tr key={s.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4">
                      <Link to={`/app/suppliers/${s.id}`} className="font-semibold text-on-surface hover:text-primary transition-colors">{s.name}</Link>
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{s.contact_person || "—"}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{s.phone || "—"}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{s.email || "—"}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        s.is_active ? "bg-secondary-container/20 text-secondary" : "bg-outline-variant/20 text-on-surface-variant"
                      }`}>{s.is_active ? "Active" : "Inactive"}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-2">
                        <Link to={`/app/suppliers/${s.id}/edit`} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">Edit</Link>
                        <button onClick={() => handleDelete(s.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </Panel>
      {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} suppliers</p>
            <div className="flex space-x-2">
              <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"><Icon name="chevron_left" /></button>
              {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                const start = Math.max(1, Math.min(page - 2, totalPages - 4)); const p = start + i;
                if (p > totalPages) return null;
                return (
                  <button key={p} onClick={() => setPage(p)}
                    className={`w-10 h-10 rounded-md flex items-center justify-center font-bold text-sm ${
                      p === page ? "bg-primary text-on-primary" : "border border-outline-variant hover:bg-surface-container-high"
                    }`}>{p}</button>
                );
              })}
              <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"><Icon name="chevron_right" /></button>
            </div>
          </div>
        )}
      </div>
  );
}

import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listDispensing, type DispenseRecord } from "../api/dispensing";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

const statusColors: Record<string, string> = {
  pending: "bg-tertiary-container/20 text-tertiary",
  in_progress: "bg-primary/10 text-primary",
  completed: "bg-secondary-container/20 text-secondary",
  cancelled: "bg-outline-variant/20 text-on-surface-variant",
};

export default function DispensingQueue() {
  const [records, setRecords] = useState<DispenseRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const limit = 20;

  const load = () => {
    setLoading(true);
    listDispensing(page, limit, status).then(res => {
      setRecords(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Dispensing Queue</h2>
          <p className="text-body-lg text-on-surface-variant">Process and track medication dispensing</p>
        </div>
      </div>

      <div className="mb-8">
        <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }}
          className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
          <option value="">All statuses</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="completed">Completed</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading dispensing records...</p></div>
          ) : records.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No dispensing records found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Patient</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Product</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Pharmacist</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Dispensed</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Controlled</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {records.map(dr => (
                  <tr key={dr.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-semibold text-on-surface">{dr.patient_name}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{dr.product_name}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{dr.pharmacist_name}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${statusColors[dr.status] || ""}`}>{dr.status === "in_progress" ? "In Progress" : dr.status}</span>
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{dr.dispensed_at ? new Date(dr.dispensed_at).toLocaleDateString() : "—"}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        dr.is_controlled ? "bg-error-container/20 text-error" : "bg-secondary-container/20 text-secondary"
                      }`}>{dr.is_controlled ? "Yes" : "No"}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <Link to={`/app/dispensing/${dr.id}`} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">View</Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} records</p>
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
    </div>
  );
}

import { useEffect, useState } from "react";
import { listSales, updateSale, getReceipt, type Sale } from "../api/pos";
import { useToast } from "../context/ToastContext";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

const statusColors: Record<string, string> = {
  active: "bg-primary/10 text-primary",
  completed: "bg-secondary-container/20 text-secondary",
  held: "bg-tertiary-container/20 text-tertiary",
  voided: "bg-outline-variant/20 text-on-surface-variant",
  refunded: "bg-error-container/20 text-error",
};

export default function SalesHistory() {
  const { showToast } = useToast();
  const [sales, setSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const limit = 20;

  const load = () => {
    setLoading(true);
    listSales(page, limit, status).then(res => {
      setSales(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / limit);

  const handleVoid = async (id: string) => {
    if (!confirm("Void this sale?")) return;
    try { await updateSale(id, { status: "voided" }); load(); }
    catch (err) { showToast(err instanceof Error ? err.message : "Void failed", "error"); }
  };

  const handleRefund = async (id: string) => {
    if (!confirm("Refund this sale?")) return;
    try { await updateSale(id, { status: "refunded" }); load(); }
    catch (err) { showToast(err instanceof Error ? err.message : "Refund failed", "error"); }
  };

  const handleReceipt = async (id: string) => {
    try {
      const receipt = await getReceipt(id);
      showToast(JSON.stringify(receipt, null, 2).slice(0, 200));
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Receipt failed", "error");
    }
  };

  return (
    <div>
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Sales History</h2>
          <p className="text-body-lg text-on-surface-variant">View and manage completed transactions</p>
        </div>
      </div>

      <div className="mb-8">
        <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }}
          className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
          <option value="">All statuses</option>
          <option value="active">Active</option>
          <option value="completed">Completed</option>
          <option value="held">Held</option>
          <option value="voided">Voided</option>
          <option value="refunded">Refunded</option>
        </select>
      </div>

      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading sales...</p></div>
          ) : sales.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No sales found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Type</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">ID</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Patient</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Total</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Paid</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Date</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {sales.map(s => (
                  <tr key={s.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4">
                      <span className="font-data-mono text-data-mono">{s.sale_type === "prescription" ? "Rx" : "OTC"}</span>
                    </td>
                    <td className="px-6 py-4 font-data-mono text-data-mono text-on-surface-variant">{s.id.slice(0, 8)}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{s.patient_name || "—"}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${statusColors[s.status] || ""}`}>{s.status}</span>
                    </td>
                    <td className="px-6 py-4 text-right font-data-mono text-data-mono">${Number(s.grand_total).toFixed(2)}</td>
                    <td className="px-6 py-4 text-right font-data-mono text-data-mono">${Number(s.paid_amount).toFixed(2)}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{s.created_at ? new Date(s.created_at).toLocaleDateString() : "—"}</td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-1">
                        <button onClick={() => handleReceipt(s.id)} className="px-2 py-1 text-xs text-primary hover:bg-primary/5 rounded transition-colors">Receipt</button>
                        {s.status === "completed" && (
                          <button onClick={() => handleRefund(s.id)} className="px-2 py-1 text-xs text-tertiary hover:bg-tertiary/5 rounded transition-colors">Refund</button>
                        )}
                        {(s.status === "active" || s.status === "completed") && (
                          <button onClick={() => handleVoid(s.id)} className="px-2 py-1 text-xs text-error hover:bg-error/5 rounded transition-colors">Void</button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} sales</p>
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

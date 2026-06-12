import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listPurchaseOrders, type PurchaseOrder } from "../api/purchases";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

const statusColors: Record<string, string> = {
  draft: "bg-outline-variant/20 text-on-surface-variant",
  approved: "bg-primary/10 text-primary",
  received: "bg-secondary-container/20 text-secondary",
  rejected: "bg-error-container/20 text-error",
};

export default function PurchaseOrders() {
  const [orders, setOrders] = useState<PurchaseOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const limit = 20;

  const load = () => {
    setLoading(true);
    listPurchaseOrders(page, limit, status).then(res => {
      setOrders(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Purchase Orders</h2>
          <p className="text-body-lg text-on-surface-variant">Track and manage procurement orders</p>
        </div>
        <Link to="/app/purchases/new" className="flex items-center px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container shadow-md transition-all">
          <Icon name="add" className="mr-2" />New Purchase Order
        </Link>
      </div>

      <div className="mb-8">
        <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }}
          className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
          <option value="">All statuses</option>
          <option value="draft">Draft</option>
          <option value="approved">Approved</option>
          <option value="received">Received</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>

      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading orders...</p></div>
          ) : orders.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No purchase orders found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">PO #</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Supplier</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Order Date</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Total</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {orders.map(po => (
                  <tr key={po.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-data-mono text-data-mono">
                      <Link to={`/app/purchases/${po.id}`} className="text-on-surface hover:text-primary transition-colors">{po.po_number}</Link>
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{po.supplier_name || "—"}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${statusColors[po.status] || ""}`}>
                        {po.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{po.order_date ? new Date(po.order_date).toLocaleDateString() : "—"}</td>
                    <td className="px-6 py-4 text-right font-data-mono text-data-mono">{po.grand_total ? `$${Number(po.grand_total).toFixed(2)}` : "—"}</td>
                    <td className="px-6 py-4 text-right">
                      <Link to={`/app/purchases/${po.id}`} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">View</Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} orders</p>
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

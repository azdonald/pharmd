import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listPurchaseOrders, type PurchaseOrder } from "../api/purchases";

export default function PurchaseOrders() {
  const [orders, setOrders] = useState<PurchaseOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);

  const load = () => {
    setLoading(true);
    listPurchaseOrders(page, 20, status).then(res => {
      setOrders(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / 20);

  const statusClass = (s: string) => {
    switch (s) {
      case "draft": return "badge badge-draft";
      case "approved": return "badge badge-approved";
      case "received": return "badge badge-received";
      case "rejected": return "badge badge-rejected";
      default: return "badge";
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Purchase Orders</h1>
        <Link to="/purchases/new" className="btn">New Purchase Order</Link>
      </div>

      <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }} style={{ marginBottom: 16 }}>
        <option value="">All statuses</option>
        <option value="draft">Draft</option>
        <option value="approved">Approved</option>
        <option value="received">Received</option>
        <option value="rejected">Rejected</option>
      </select>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>PO #</th><th>Supplier</th><th>Status</th><th>Order Date</th><th>Total</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {orders.map(po => (
                <tr key={po.id}>
                  <td><Link to={`/purchases/${po.id}`}>{po.po_number}</Link></td>
                  <td>{po.supplier_name || "—"}</td>
                  <td><span className={statusClass(po.status)}>{po.status}</span></td>
                  <td>{po.order_date ? new Date(po.order_date).toLocaleDateString() : "—"}</td>
                  <td>{po.grand_total ? `$${Number(po.grand_total).toFixed(2)}` : "—"}</td>
                  <td>
                    <Link to={`/purchases/${po.id}`}>View</Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {totalPages > 1 && (
            <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
              {page > 1 && <button onClick={() => setPage(page - 1)}>Previous</button>}
              <span style={{ padding: "8px 0" }}>Page {page} of {totalPages}</span>
              {page < totalPages && <button onClick={() => setPage(page + 1)}>Next</button>}
            </div>
          )}
        </>
      )}
    </div>
  );
}

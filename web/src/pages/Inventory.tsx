import { useEffect, useState } from "react";
import { listStock, listAlerts, listExpiring, type StockItem, type StockAlert, type ExpiringBatch } from "../api/inventory";
import { listLocations, type Location } from "../api/locations";

export default function Inventory() {
  const [locations, setLocations] = useState<Location[]>([]);
  const [locationId, setLocationId] = useState("");
  const [stock, setStock] = useState<StockItem[]>([]);
  const [alerts, setAlerts] = useState<StockAlert[]>([]);
  const [expiring, setExpiring] = useState<ExpiringBatch[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    listLocations().then(res => {
      setLocations(res.data);
      if (res.data.length > 0) setLocationId(res.data[0].id);
    }).catch(console.error);
  }, []);

  const load = () => {
    if (!locationId) return;
    setLoading(true);
    Promise.all([
      listStock(locationId, 1, 200, search),
      listAlerts(locationId),
      listExpiring(locationId, 30),
    ]).then(([s, a, e]) => {
      setStock(s.data);
      setAlerts(a);
      setExpiring(e);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [locationId, search]);

  const totalStock = stock.reduce((sum, item) => sum + item.total_quantity, 0);
  const lowStockCount = alerts.length;
  const expiringCount = expiring.length;

  return (
    <div>
      <div className="page-header">
        <h1>Inventory</h1>
        <div style={{ display: "flex", gap: 8 }}>
          <select value={locationId} onChange={e => setLocationId(e.target.value)}>
            {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
          </select>
          <input value={search} onChange={e => setSearch(e.target.value)} placeholder="Search..." />
        </div>
      </div>

      <div style={{ display: "flex", gap: 16, marginBottom: 20 }}>
        <div style={{ flex: 1, background: "#fff", padding: 16, borderRadius: 4, boxShadow: "0 1px 3px rgba(0,0,0,0.1)" }}>
          <strong>Total Stock</strong>
          <p style={{ fontSize: 24, margin: "8px 0 0" }}>{totalStock}</p>
        </div>
        <div style={{ flex: 1, background: "#fff", padding: 16, borderRadius: 4, boxShadow: "0 1px 3px rgba(0,0,0,0.1)" }}>
          <strong>Low Stock</strong>
          <p style={{ fontSize: 24, margin: "8px 0 0", color: lowStockCount > 0 ? "#d32f2f" : "inherit" }}>{lowStockCount}</p>
        </div>
        <div style={{ flex: 1, background: "#fff", padding: 16, borderRadius: 4, boxShadow: "0 1px 3px rgba(0,0,0,0.1)" }}>
          <strong>Expiring Soon</strong>
          <p style={{ fontSize: 24, margin: "8px 0 0", color: expiringCount > 0 ? "#e65100" : "inherit" }}>{expiringCount}</p>
        </div>
      </div>

      {loading ? <p>Loading...</p> : (
        <>
          <h2 style={{ fontSize: 18, marginBottom: 12 }}>Stock by Product</h2>
          <table>
            <thead>
              <tr>
                <th>Product</th><th>Brand</th><th>Generic</th><th>Classification</th>
                <th>Total Qty</th><th>Reorder Level</th><th>Status</th><th>Batches</th>
              </tr>
            </thead>
            <tbody>
              {stock.map(item => (
                <tr key={item.product_id}>
                  <td>{item.product_name}</td>
                  <td>{item.brand_name || "—"}</td>
                  <td>{item.generic_name || "—"}</td>
                  <td>{item.classification || "—"}</td>
                  <td>{item.total_quantity}</td>
                  <td>{item.reorder_level}</td>
                  <td style={{ color: item.total_quantity < item.reorder_level ? "#d32f2f" : "#2e7d32" }}>
                    {item.total_quantity < item.reorder_level ? "Low" : "OK"}
                  </td>
                  <td>
                    <details>
                      <summary>{item.batches.length} batch(es)</summary>
                      <table style={{ fontSize: 12, marginTop: 4 }}>
                        <thead>
                          <tr><th>Batch</th><th>Qty</th><th>Cost</th><th>Price</th><th>Expiry</th></tr>
                        </thead>
                        <tbody>
                          {item.batches.map(b => (
                            <tr key={b.id}>
                              <td><code>{b.batch_number || "—"}</code></td>
                              <td>{b.quantity}</td>
                              <td>{b.unit_cost?.toFixed(2)}</td>
                              <td>{b.selling_price?.toFixed(2)}</td>
                              <td style={{ color: b.expiry_date && new Date(b.expiry_date) < new Date(Date.now() + 30*86400000) ? "#e65100" : "inherit" }}>
                                {b.expiry_date || "—"}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </details>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      {alerts.length > 0 && (
        <>
          <h2 style={{ fontSize: 18, marginTop: 24, marginBottom: 12, color: "#d32f2f" }}>Low Stock Alerts</h2>
          <table>
            <thead><tr><th>Product</th><th>Brand</th><th>Total Qty</th><th>Reorder Level</th></tr></thead>
            <tbody>
              {alerts.map(a => (
                <tr key={a.product_id}>
                  <td>{a.product_name}</td>
                  <td>{a.brand_name || "—"}</td>
                  <td style={{ color: "#d32f2f", fontWeight: 600 }}>{a.total_quantity}</td>
                  <td>{a.reorder_level}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      {expiring.length > 0 && (
        <>
          <h2 style={{ fontSize: 18, marginTop: 24, marginBottom: 12, color: "#e65100" }}>Expiring Soon</h2>
          <table>
            <thead><tr><th>Product</th><th>Batch</th><th>Remaining</th><th>Expiry</th><th>Days Left</th></tr></thead>
            <tbody>
              {expiring.map(e => (
                <tr key={e.id}>
                  <td>{e.product_name}</td>
                  <td><code>{e.batch_number || "—"}</code></td>
                  <td>{e.remaining_qty}</td>
                  <td>{e.expiry_date}</td>
                  <td style={{ color: e.days_until_expiry <= 7 ? "#d32f2f" : "#e65100" }}>{e.days_until_expiry} days</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  );
}

import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getPurchaseOrder, createPurchaseOrder, approvePurchaseOrder, rejectPurchaseOrder, receiveGoods, type PurchaseOrder, type POItemInput } from "../api/purchases";
import { listLocations, type Location } from "../api/locations";
import { listSuppliers, type Supplier } from "../api/suppliers";
import { listProducts, type Product } from "../api/products";
import { useToast } from "../context/ToastContext";

export default function PurchaseOrderDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [po, setPo] = useState<PurchaseOrder | null>(null);
  const [loading, setLoading] = useState(true);
  const { showToast } = useToast();
  const [locations, setLocations] = useState<Location[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [products, setProducts] = useState<Product[]>([]);

  // receive modal
  const [showReceive, setShowReceive] = useState(false);
  const [receiveItems, setReceiveItems] = useState<{ itemId: string; received: number }[]>([]);
  const [receiveNote, setReceiveNote] = useState("");

  const load = () => {
    if (!id) return;
    setLoading(true);
    getPurchaseOrder(id)
      .then(poRes => Promise.all([poRes, listLocations(1, 100), listSuppliers(1, 200), listProducts(1, 200)]))
      .then(([poRes, locRes, supRes, prodRes]) => {
        setPo(poRes);
        setLocations(locRes.data);
        setSuppliers(supRes.data);
        setProducts(prodRes.data);
      })
      .catch(err => { console.error(err); showToast("Failed to load purchase order", "error"); })
      .finally(() => setLoading(false));
  };

  useEffect(load, [id]);

  const handleApprove = async () => {
    if (!id) return;
    try {
      const updated = await approvePurchaseOrder(id);
      setPo(updated);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Approve failed", "error");
    }
  };

  const handleReject = async () => {
    if (!id) return;
    try {
      const updated = await rejectPurchaseOrder(id);
      setPo(updated);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Reject failed", "error");
    }
  };

  const openReceive = () => {
    if (!po?.items) return;
    setReceiveItems(po.items.map(i => ({ itemId: i.id, received: 0 })));
    setShowReceive(true);
  };

  const handleReceive = async () => {
    if (!id) return;
    try {
      const items = receiveItems.filter(i => i.received > 0).map(i => ({ item_id: i.itemId, quantity_received: i.received }));
      if (items.length === 0) { showToast("At least one item with quantity > 0", "error"); return; }
      const updated = await receiveGoods(id, { items, notes: receiveNote || undefined });
      setPo(updated);
      setShowReceive(false);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Receive failed", "error");
    }
  };

  const getLocationName = (id: string) => locations.find(l => l.id === id)?.name || id;
  const getSupplierName = (id: string) => suppliers.find(s => s.id === id)?.name || id;
  const getProductName = (id: string) => products.find(p => p.id === id)?.name || id;

  if (loading) return <p>Loading...</p>;
  if (!po) return <p>Purchase order not found</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Purchase Order: {po.po_number}</h1>
        <button onClick={() => navigate("/purchases")}>Back</button>
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 24 }}>
        <div>
          <p><strong>Supplier:</strong> {getSupplierName(po.supplier_id)}</p>
          <p><strong>Location:</strong> {getLocationName(po.location_id)}</p>
          <p><strong>Status:</strong> <span className={`badge badge-${po.status}`}>{po.status}</span></p>
        </div>
        <div>
          <p><strong>Order Date:</strong> {po.order_date ? new Date(po.order_date).toLocaleDateString() : "—"}</p>
          <p><strong>Expected Date:</strong> {po.expected_date || "—"}</p>
          <p><strong>Notes:</strong> {po.notes || "—"}</p>
        </div>
      </div>

      <h3>Items</h3>
      <table>
        <thead>
          <tr><th>Product</th><th>Ordered</th><th>Received</th><th>Unit Cost</th><th>Line Total</th></tr>
        </thead>
        <tbody>
          {(po.items || []).map(item => (
            <tr key={item.id}>
              <td>{getProductName(item.product_id)}</td>
              <td>{item.quantity_ordered}</td>
              <td>{item.quantity_received}</td>
              <td>${Number(item.unit_cost).toFixed(2)}</td>
              <td>${Number(item.line_total).toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
        <tfoot>
          <tr>
            <td colSpan={4} style={{ textAlign: "right" }}><strong>Subtotal:</strong></td>
            <td>${Number(po.subtotal).toFixed(2)}</td>
          </tr>
          <tr>
            <td colSpan={4} style={{ textAlign: "right" }}><strong>Tax:</strong></td>
            <td>${Number(po.tax_total).toFixed(2)}</td>
          </tr>
          <tr>
            <td colSpan={4} style={{ textAlign: "right" }}><strong>Grand Total:</strong></td>
            <td>${Number(po.grand_total).toFixed(2)}</td>
          </tr>
        </tfoot>
      </table>

      <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
        {po.status === "draft" && <button onClick={handleApprove}>Approve</button>}
        {po.status === "draft" && <button onClick={handleReject}>Reject</button>}
        {po.status === "approved" && <button onClick={openReceive}>Receive Goods</button>}
      </div>

      {showReceive && (
        <div className="modal">
          <div className="modal-content">
            <h2>Receive Goods</h2>
            <table>
              <thead>
                <tr><th>Product</th><th>Ordered</th><th>Previously Received</th><th>Qty Receiving</th></tr>
              </thead>
              <tbody>
                {receiveItems.map((ri, idx) => {
                  const item = (po.items || [])[idx];
                  return (
                    <tr key={ri.itemId}>
                      <td>{getProductName(item.product_id)}</td>
                      <td>{item.quantity_ordered}</td>
                      <td>{item.quantity_received}</td>
                      <td>
                        <input type="number" min={0} max={item.quantity_ordered - (item.quantity_received || 0)} value={ri.received} onChange={e => {
                          const updated = [...receiveItems];
                          updated[idx] = { ...updated[idx], received: Number(e.target.value) };
                          setReceiveItems(updated);
                        }} />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
            <div style={{ marginTop: 8 }}>
              <label>Notes:</label>
              <textarea value={receiveNote} onChange={e => setReceiveNote(e.target.value)} rows={2} style={{ width: "100%" }} />
            </div>
            <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
              <button onClick={handleReceive}>Confirm Receipt</button>
              <button onClick={() => setShowReceive(false)}>Cancel</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Inline PO creation form used by PurchaseOrderForm
export function PurchaseOrderForm() {
  const navigate = useNavigate();
  const { showToast } = useToast();
  const [locations, setLocations] = useState<Location[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [products, setProducts] = useState<Product[]>([]);

  const [supplierId, setSupplierId] = useState("");
  const [locationId, setLocationId] = useState("");
  const [expectedDate, setExpectedDate] = useState("");
  const [notes, setNotes] = useState("");
  const [items, setItems] = useState<POItemInput[]>([]);

  useEffect(() => {
    Promise.all([
      listLocations(1, 100),
      listSuppliers(1, 200),
      listProducts(1, 200),
    ]).then(([locRes, supRes, prodRes]) => {
      setLocations(locRes.data);
      setSuppliers(supRes.data);
      setProducts(prodRes.data);
    }).catch(console.error);
  }, []);

  const addItem = () => {
    setItems([...items, { product_id: "", quantity_ordered: 1, unit_cost: 0 }]);
  };

  const updateItem = (idx: number, field: keyof POItemInput, value: string | number) => {
    const updated = [...items];
    (updated[idx] as any)[field] = value;
    setItems(updated);
  };

  const removeItem = (idx: number) => {
    setItems(items.filter((_, i) => i !== idx));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!supplierId || !locationId || items.length === 0) {
      showToast("Please fill all required fields and add at least one item", "error");
      return;
    }
    if (items.some(i => !i.product_id)) {
      showToast("Please select a product for each item", "error");
      return;
    }
    try {
      const created = await createPurchaseOrder({
        supplier_id: supplierId,
        location_id: locationId,
        expected_date: expectedDate || undefined,
        notes: notes || undefined,
        items,
      });
      navigate(`/purchases/${created.id}`);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Create failed", "error");
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>New Purchase Order</h1>
      </div>
      <div className="form-page">
      <form onSubmit={handleSubmit}>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 16 }}>
          <div>
            <label>Supplier *</label>
            <select value={supplierId} onChange={e => setSupplierId(e.target.value)} required>
              <option value="">Select supplier</option>
              {suppliers.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
            </select>
          </div>
          <div>
            <label>Location *</label>
            <select value={locationId} onChange={e => setLocationId(e.target.value)} required>
              <option value="">Select location</option>
              {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
            </select>
          </div>
          <div>
            <label>Expected Date</label>
            <input type="date" value={expectedDate} onChange={e => setExpectedDate(e.target.value)} />
          </div>
          <div>
            <label>Notes</label>
            <textarea value={notes} onChange={e => setNotes(e.target.value)} rows={2} />
          </div>
        </div>

        <h3>Items</h3>
        <table>
          <thead>
            <tr><th>Product</th><th>Qty</th><th>Unit Cost</th><th>Line Total</th><th></th></tr>
          </thead>
          <tbody>
              {items.map((item, idx) => {
              const lineTotal = item.quantity_ordered * item.unit_cost;
              return (
                <tr key={idx}>
                  <td>
                    <select value={item.product_id} onChange={e => updateItem(idx, "product_id", e.target.value)} required>
                      <option value="">Select product</option>
                      {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
                    </select>
                  </td>
                  <td><input type="number" min={1} value={item.quantity_ordered} onChange={e => updateItem(idx, "quantity_ordered", Number(e.target.value))} /></td>
                  <td><input type="number" min={0} step={0.01} value={item.unit_cost} onChange={e => updateItem(idx, "unit_cost", Number(e.target.value))} /></td>
                  <td>${lineTotal.toFixed(2)}</td>
                  <td><button type="button" onClick={() => removeItem(idx)}>Remove</button></td>
                </tr>
              );
            })}
          </tbody>
        </table>
        <button type="button" onClick={addItem} style={{ marginTop: 8 }}>Add Item</button>

        <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
          <button type="submit">Create Purchase Order</button>
          <button type="button" onClick={() => navigate("/purchases")}>Cancel</button>
        </div>
      </form>
      </div>
    </div>
  );
}

import { useEffect, useState } from "react";
import { listProducts, type Product } from "../api/products";
import { listLocations, type Location } from "../api/locations";
import { listPatients, type Patient } from "../api/patients";
import { createSale, recordPayment, getReceipt, getDailySummary, closeDay, type Sale } from "../api/pos";

export default function POS() {
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [patients, setPatients] = useState<Patient[]>([]);
  const [cart, setCart] = useState<{ product: Product; qty: number; price: number; discount: number }[]>([]);
  const [selectedLoc, setSelectedLoc] = useState("");
  const [selectedPatient, setSelectedPatient] = useState("");
  const [productSearch, setProductSearch] = useState("");
  const [saleNotes, setSaleNotes] = useState("");
  const [showPayment, setShowPayment] = useState(false);
  const [cashAmount, setCashAmount] = useState("");
  const [cardAmount, setCardAmount] = useState("");
  const [mobileAmount, setMobileAmount] = useState("");
  const [mobileRef, setMobileRef] = useState("");
  const [lastSale, setLastSale] = useState<Sale | null>(null);
  const [summary, setSummary] = useState<any>(null);
  const [summaryDate, setSummaryDate] = useState(new Date().toISOString().split("T")[0]);

  useEffect(() => {
    Promise.all([
      listLocations(1, 100),
      listProducts(1, 500),
      listPatients(1, 200),
    ]).then(([locRes, prodRes, patRes]) => {
      setLocations(locRes.data);
      setProducts(prodRes.data);
      setPatients(patRes.data);
    }).catch(console.error);
  }, []);

  const filteredProducts = products.filter(p =>
    !productSearch || p.name.toLowerCase().includes(productSearch.toLowerCase())
  );

  const addToCart = (product: Product) => {
    const existing = cart.find(c => c.product.id === product.id);
    if (existing) {
      setCart(cart.map(c => c.product.id === product.id ? { ...c, qty: c.qty + 1 } : c));
    } else {
      setCart([...cart, { product, qty: 1, price: 0, discount: 0 }]);
    }
  };

  const updateCart = (idx: number, field: string, value: any) => {
    const updated = [...cart];
    (updated[idx] as any)[field] = value;
    setCart(updated);
  };

  const removeFromCart = (idx: number) => {
    setCart(cart.filter((_, i) => i !== idx));
  };

  const subtotal = cart.reduce((sum, c) => sum + c.qty * c.price, 0);
  const discountTotal = cart.reduce((sum, c) => sum + c.discount, 0);
  const grandTotal = subtotal - discountTotal;

  const handleCheckout = async () => {
    if (!selectedLoc || cart.length === 0) { alert("Select a location and add items"); return; }
    if (cart.some(c => c.price <= 0)) { alert("Set price for all items"); return; }
    try {
      const sale = await createSale({
        location_id: selectedLoc,
        patient_id: selectedPatient || undefined,
        sale_type: "otc",
        notes: saleNotes || undefined,
        items: cart.map(c => ({ product_id: c.product.id, quantity: c.qty, unit_price: c.price, discount: c.discount || undefined })),
      });
      setLastSale(sale);
      setShowPayment(true);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Checkout failed");
    }
  };

  const handlePayment = async () => {
    if (!lastSale) return;
    const payments: { method: string; amount: number; reference?: string }[] = [];
    if (Number(cashAmount) > 0) payments.push({ method: "cash", amount: Number(cashAmount) });
    if (Number(cardAmount) > 0) payments.push({ method: "card", amount: Number(cardAmount) });
    if (Number(mobileAmount) > 0) payments.push({ method: "mobile_money", amount: Number(mobileAmount), reference: mobileRef || undefined });
    if (payments.length === 0) { alert("Enter at least one payment"); return; }
    try {
      await recordPayment(lastSale.id, payments);
      setCart([]);
      setShowPayment(false);
      setLastSale(null);
      setCashAmount("");
      setCardAmount("");
      setMobileAmount("");
      setMobileRef("");
      setSaleNotes("");
      alert("Sale completed!");
    } catch (err) {
      alert(err instanceof Error ? err.message : "Payment failed");
    }
  };

  const handleReceipt = async () => {
    if (!lastSale) return;
    const receipt = await getReceipt(lastSale.id);
    alert(JSON.stringify(receipt, null, 2));
  };

  const loadSummary = async () => {
    if (!selectedLoc) { alert("Select a location"); return; }
    try {
      const s = await getDailySummary(summaryDate, selectedLoc);
      setSummary(s);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to load summary");
    }
  };

  const handleCloseDay = async () => {
    if (!selectedLoc || !summaryDate) { alert("Select location and date"); return; }
    if (!confirm("Close day for " + summaryDate + "?")) return;
    try {
      const s = await closeDay(summaryDate, selectedLoc);
      setSummary(s);
      alert("Day closed!");
    } catch (err) {
      alert(err instanceof Error ? err.message : "Close day failed");
    }
  };

  return (
    <div style={{ display: "flex", gap: 16, height: "calc(100vh - 100px)" }}>
      {/* Product Panel */}
      <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
        <div className="page-header">
          <h1>Point of Sale</h1>
        </div>
        <div style={{ display: "flex", gap: 8, marginBottom: 8 }}>
          <select value={selectedLoc} onChange={e => setSelectedLoc(e.target.value)} style={{ flex: 1 }}>
            <option value="">Select location</option>
            {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
          </select>
          <select value={selectedPatient} onChange={e => setSelectedPatient(e.target.value)} style={{ flex: 1 }}>
            <option value="">Walk-in (no patient)</option>
            {patients.map(p => <option key={p.id} value={p.id}>{p.first_name} {p.last_name}</option>)}
          </select>
        </div>
        <input value={productSearch} onChange={e => setProductSearch(e.target.value)} placeholder="Search products..." style={{ marginBottom: 8 }} />
        <div style={{ flex: 1, overflow: "auto", border: "1px solid #ddd", borderRadius: 4 }}>
          {filteredProducts.map(p => (
            <div key={p.id} onClick={() => addToCart(p)} style={{ padding: "8px 12px", cursor: "pointer", borderBottom: "1px solid #eee", display: "flex", justifyContent: "space-between" }}>
              <span>{p.name}</span>
              <span style={{ color: "#666" }}>{p.form || ""}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Cart Panel */}
      <div style={{ width: 400, display: "flex", flexDirection: "column" }}>
        <h3>Cart ({cart.length} items)</h3>
        <div style={{ flex: 1, overflow: "auto", border: "1px solid #ddd", borderRadius: 4, marginBottom: 8 }}>
          {cart.map((c, idx) => (
            <div key={c.product.id} style={{ padding: 8, borderBottom: "1px solid #eee" }}>
              <div style={{ display: "flex", justifyContent: "space-between" }}>
                <strong>{c.product.name}</strong>
                <button onClick={() => removeFromCart(idx)}>x</button>
              </div>
              <div style={{ display: "flex", gap: 4, marginTop: 4 }}>
                <input type="number" min={1} value={c.qty} onChange={e => updateCart(idx, "qty", Number(e.target.value))} style={{ width: 60 }} />
                <input type="number" min={0} step={0.01} value={c.price} onChange={e => updateCart(idx, "price", Number(e.target.value))} placeholder="Price" style={{ flex: 1 }} />
                <input type="number" min={0} step={0.01} value={c.discount} onChange={e => updateCart(idx, "discount", Number(e.target.value))} placeholder="Disc" style={{ width: 60 }} />
                <span>${(c.qty * c.price - c.discount).toFixed(2)}</span>
              </div>
            </div>
          ))}
        </div>
        <div style={{ borderTop: "2px solid #333", padding: "8px 0" }}>
          <p><strong>Subtotal:</strong> ${subtotal.toFixed(2)}</p>
          <p><strong>Discounts:</strong> -${discountTotal.toFixed(2)}</p>
          <p><strong>Grand Total:</strong> ${grandTotal.toFixed(2)}</p>
        </div>
        <textarea value={saleNotes} onChange={e => setSaleNotes(e.target.value)} placeholder="Notes..." rows={2} style={{ marginBottom: 8 }} />
        <button onClick={handleCheckout} disabled={cart.length === 0} style={{ padding: 12, fontSize: 16 }}>Checkout (${grandTotal.toFixed(2)})</button>
      </div>

      {/* Payment Modal */}
      {showPayment && lastSale && (
        <div className="modal">
          <div className="modal-content">
            <h2>Payment — ${lastSale.grand_total?.toFixed(2)}</h2>
            <div style={{ marginBottom: 8 }}>
              <label>Cash</label>
              <input type="number" min={0} step={0.01} value={cashAmount} onChange={e => setCashAmount(e.target.value)} />
            </div>
            <div style={{ marginBottom: 8 }}>
              <label>Card</label>
              <input type="number" min={0} step={0.01} value={cardAmount} onChange={e => setCardAmount(e.target.value)} />
            </div>
            <div style={{ marginBottom: 8 }}>
              <label>Mobile Money</label>
              <input type="number" min={0} step={0.01} value={mobileAmount} onChange={e => setMobileAmount(e.target.value)} />
              <input value={mobileRef} onChange={e => setMobileRef(e.target.value)} placeholder="Reference" style={{ marginTop: 4 }} />
            </div>
            <div style={{ display: "flex", gap: 8 }}>
              <button onClick={handlePayment}>Complete Payment</button>
              <button onClick={() => setShowPayment(false)}>Cancel</button>
            </div>
            {lastSale.id && <button onClick={handleReceipt} style={{ marginTop: 8 }}>Print Receipt</button>}
          </div>
        </div>
      )}

      {/* Summary Panel */}
      <div style={{ width: 300, borderLeft: "1px solid #ddd", paddingLeft: 16 }}>
        <h3>Reports</h3>
        <div style={{ marginBottom: 8 }}>
          <label>Date</label>
          <input type="date" value={summaryDate} onChange={e => setSummaryDate(e.target.value)} style={{ width: "100%" }} />
        </div>
        <button onClick={loadSummary} style={{ width: "100%", marginBottom: 8 }}>X-Report (Summary)</button>
        <button onClick={handleCloseDay} style={{ width: "100%" }}>Z-Report (Close Day)</button>
        {summary && (
          <div style={{ marginTop: 16, padding: 12, border: "1px solid #ddd", borderRadius: 4 }}>
            <p><strong>Date:</strong> {summary.date}</p>
            <p><strong>Sales:</strong> {summary.total_sales}</p>
            <p><strong>Revenue:</strong> ${Number(summary.total_revenue).toFixed(2)}</p>
            <p><strong>Tax:</strong> ${Number(summary.total_tax).toFixed(2)}</p>
            <p><strong>Discounts:</strong> ${Number(summary.total_discounts).toFixed(2)}</p>
            <p><strong>Closed:</strong> {summary.is_closed ? "Yes" : "No"}</p>
            {summary.by_method && Object.entries(summary.by_method).map(([method, amt]) => (
              <p key={method}>{method}: ${Number(amt).toFixed(2)}</p>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

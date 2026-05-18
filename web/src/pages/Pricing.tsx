import { useEffect, useState } from "react";
import { listPrices, upsertPrice, deletePrice, type ProductPrice } from "../api/pricing";
import { listLocations, type Location } from "../api/locations";
import { listProducts, type Product } from "../api/products";

export default function Pricing() {
  const [prices, setPrices] = useState<ProductPrice[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [productFilter, setProductFilter] = useState("");
  const [locationFilter, setLocationFilter] = useState("");
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);

  const [showForm, setShowForm] = useState(false);
  const [formProduct, setFormProduct] = useState("");
  const [formLocation, setFormLocation] = useState("");
  const [formSelling, setFormSelling] = useState("");

  const load = () => {
    setLoading(true);
    Promise.all([
      listPrices(productFilter, locationFilter, page, 20),
      listLocations(1, 100),
      listProducts(1, 200),
    ]).then(([priceRes, locRes, prodRes]) => {
      setPrices(priceRes.data);
      setTotal(priceRes.total);
      setLocations(locRes.data);
      setProducts(prodRes.data);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, productFilter, locationFilter]);

  const totalPages = Math.ceil(total / 20);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formProduct || !formLocation || !formSelling) return;
    try {
      await upsertPrice({ product_id: formProduct, location_id: formLocation, selling_price: Number(formSelling) });
      setShowForm(false);
      setFormProduct("");
      setFormLocation("");
      setFormSelling("");
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to save price");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this price?")) return;
    try {
      await deletePrice(id);
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Product Pricing</h1>
        <button onClick={() => setShowForm(!showForm)}>{showForm ? "Cancel" : "Add Price"}</button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 16, padding: 16, border: "1px solid #ddd", borderRadius: 4 }}>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr auto", gap: 8, alignItems: "end" }}>
            <div>
              <label>Product *</label>
              <select value={formProduct} onChange={e => setFormProduct(e.target.value)} required>
                <option value="">Select product</option>
                {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
              </select>
            </div>
            <div>
              <label>Location *</label>
              <select value={formLocation} onChange={e => setFormLocation(e.target.value)} required>
                <option value="">Select location</option>
                {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
              </select>
            </div>
            <div>
              <label>Selling Price *</label>
              <input type="number" min={0} step={0.01} value={formSelling} onChange={e => setFormSelling(e.target.value)} required />
            </div>
            <button type="submit">Save</button>
          </div>
        </form>
      )}

      <div style={{ display: "flex", gap: 8, marginBottom: 16 }}>
        <select value={productFilter} onChange={e => { setProductFilter(e.target.value); setPage(1); }}>
          <option value="">All products</option>
          {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
        </select>
        <select value={locationFilter} onChange={e => { setLocationFilter(e.target.value); setPage(1); }}>
          <option value="">All locations</option>
          {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
        </select>
      </div>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Product</th><th>Location</th><th>Selling Price</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {prices.map(p => (
                <tr key={p.id}>
                  <td>{p.product_name}</td>
                  <td>{p.location_name}</td>
                  <td>${Number(p.selling_price).toFixed(2)}</td>
                  <td><button onClick={() => handleDelete(p.id)}>Delete</button></td>
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

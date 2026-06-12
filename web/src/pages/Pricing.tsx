import { useEffect, useState } from "react";
import { listPrices, upsertPrice, deletePrice, type ProductPrice } from "../api/pricing";
import { listLocations, type Location } from "../api/locations";
import { listProducts, type Product } from "../api/products";
import { useToast } from "../context/ToastContext";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

export default function Pricing() {
  const { showToast } = useToast();
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
  const limit = 20;

  const load = () => {
    setLoading(true);
    Promise.all([
      listPrices(productFilter, locationFilter, page, limit),
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

  const totalPages = Math.ceil(total / limit);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formProduct || !formLocation || !formSelling) return;
    try {
      await upsertPrice({ product_id: formProduct, location_id: formLocation, selling_price: Number(formSelling) });
      setShowForm(false);
      setFormProduct(""); setFormLocation(""); setFormSelling("");
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed to save price", "error");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this price?")) return;
    try {
      await deletePrice(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  return (
    <div>
      <div className="flex justify-between items-end mb-8">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Product Pricing</h2>
          <p className="text-body-lg text-on-surface-variant">Manage product selling prices per location</p>
        </div>
        <button onClick={() => setShowForm(!showForm)}
          className={`flex items-center px-4 py-2 font-semibold rounded-lg transition-all ${
            showForm
              ? "border border-outline-variant text-on-surface hover:bg-surface-container-high"
              : "bg-primary text-on-primary shadow-md hover:bg-primary-container"
          }`}>
          <Icon name={showForm ? "close" : "add"} className="mr-2" />
          {showForm ? "Cancel" : "Add Price"}
        </button>
      </div>

      {/* Inline form */}
      {showForm && (
        <form onSubmit={handleSubmit} className="mb-8 p-4 rounded-xl border border-outline-variant bg-surface-container-lowest">
          <h3 className="font-semibold text-on-surface mb-3">New Price Entry</h3>
          <div className="flex gap-3 items-end">
            <div className="flex-1">
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Product *</label>
              <select value={formProduct} onChange={e => setFormProduct(e.target.value)} required
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
                <option value="">Select product</option>
                {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
              </select>
            </div>
            <div className="flex-1">
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Location *</label>
              <select value={formLocation} onChange={e => setFormLocation(e.target.value)} required
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
                <option value="">Select location</option>
                {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Selling Price *</label>
              <input type="number" min={0} step={0.01} value={formSelling} onChange={e => setFormSelling(e.target.value)} required
                className="w-32 rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <button type="submit"
              className="px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container transition-all">Save</button>
          </div>
        </form>
      )}

      {/* Filters */}
      <div className="mb-8 flex gap-3">
        <select value={productFilter} onChange={e => { setProductFilter(e.target.value); setPage(1); }}
          className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
          <option value="">All products</option>
          {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
        </select>
        <select value={locationFilter} onChange={e => { setLocationFilter(e.target.value); setPage(1); }}
          className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
          <option value="">All locations</option>
          {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
        </select>
      </div>

      {/* Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading prices...</p></div>
          ) : prices.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No prices found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Product</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Location</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Selling Price</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {prices.map(p => (
                  <tr key={p.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-semibold text-on-surface">{p.product_name}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.location_name}</td>
                    <td className="px-6 py-4 text-right font-data-mono text-data-mono">${Number(p.selling_price).toFixed(2)}</td>
                    <td className="px-6 py-4 text-right">
                      <button onClick={() => handleDelete(p.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} prices</p>
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

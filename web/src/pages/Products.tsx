import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { listProducts, listCategories, deleteProduct, type Product, type ProductCategory } from "../api/products";
import { useToast } from "../context/ToastContext";
import { Icon, PageHeader, Panel } from "../components/AdminComponents";


export default function Products() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState(searchParams.get("query") || "");
  const { showToast } = useToast();

  const page = parseInt(searchParams.get("page") || "1");
  const query = searchParams.get("query") || "";
  const categoryId = searchParams.get("category_id") || "";
  const limit = 20;

  useEffect(() => {
    listCategories().then(setCategories).catch(console.error);
  }, []);

  const load = () => {
    setLoading(true);
    listProducts(page, limit, query, categoryId)
      .then(res => {
        setProducts(res.data);
        setTotal(res.total);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  };

  useEffect(load, [page, query, categoryId]);

  const catMap = Object.fromEntries(categories.map(c => [c.id, c.name]));

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ query: search, page: "1" });
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this product?")) return;
    try {
      await deleteProduct(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      <PageHeader
        title="Products"
        description="Manage medication and product catalog"
        actions={
          <Link
            to="/app/products/new"
            className="btn-sky-action"
          >
            <Icon name="add" className="mr-2" />
            New Product
          </Link>
        }
      />

      {/* Search & Filters */}
      <form onSubmit={handleSearch} className="mb-8 flex gap-3 items-end">
        <div className="relative flex-1 max-w-md">
          <Icon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant" />
          <input
            value={search}
            onChange={e => setSearch(e.target.value)}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest py-2 pl-10 pr-4 text-sm outline-none focus:ring-2 focus:ring-primary"
            placeholder="Search by name, brand, barcode..."
            type="text"
          />
        </div>
        <select
          value={categoryId}
          onChange={e => setSearchParams({ query, category_id: e.target.value, page: "1" })}
          className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none"
        >
          <option value="">All Categories</option>
          {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <button
          type="submit"
          className="btn-sky-action"
        >
          Search
        </button>
      </form>

      {/* Table */}
      <Panel>
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center">
              <p className="text-on-surface-variant">Loading products...</p>
            </div>
          ) : products.length === 0 ? (
            <div className="p-12 text-center">
              <p className="text-on-surface-variant">No products found</p>
            </div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Brand</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Generic</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Category</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Strength</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Form</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {products.map(p => (
                  <tr key={p.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4">
                      <Link to={`/app/products/${p.id}`} className="font-semibold text-on-surface hover:text-primary transition-colors">
                        {p.name}
                      </Link>
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.brand_name || "—"}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.generic_name || "—"}</td>
                    <td className="px-6 py-4">
                      <span className="text-body-md text-on-surface-variant">{catMap[p.category_id] || "—"}</span>
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.strength || "—"}</td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">{p.form || "—"}</td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        p.is_active
                          ? "bg-secondary-container/20 text-secondary"
                          : "bg-outline-variant/20 text-on-surface-variant"
                      }`}>
                        {p.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-2">
                        <Link to={`/app/products/${p.id}/edit`} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">Edit</Link>
                        <button onClick={() => handleDelete(p.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {/* Pagination */}
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">
              Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} products
            </p>
            <div className="flex space-x-2">
              <button onClick={() => setSearchParams({ query, category_id: categoryId, page: String(page - 1) })}
                disabled={page === 1}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed">
                <Icon name="chevron_left" />
              </button>
              {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                const start = Math.max(1, Math.min(page - 2, totalPages - 4));
                const p = start + i;
                if (p > totalPages) return null;
                return (
                  <button key={p} onClick={() => setSearchParams({ query, category_id: categoryId, page: String(p) })}
                    className={`w-10 h-10 rounded-md flex items-center justify-center font-bold text-sm ${
                      p === page ? "bg-primary text-on-primary" : "border border-outline-variant hover:bg-surface-container-high"
                    }`}>{p}</button>
                );
              })}
              <button onClick={() => setSearchParams({ query, category_id: categoryId, page: String(page + 1) })}
                disabled={page === totalPages}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed">
                <Icon name="chevron_right" />
              </button>
            </div>
          </div>
        )}
      </Panel>
      </div>
  );
}

python3 experiments/boutique/populate.py \
    --users 10000 \
    --products 10000 \
    --product_size 10000 \
    --frontend $(kip frontend) \
    --product_catalog $(kip product_catalog) \
    --currency $(kip currency)

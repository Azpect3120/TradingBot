for (( i=1; i<=1000; i++ )); do
  # curl -s "https://query1.finance.yahoo.com/v8/finance/chart/amd"
  go run ./cmd/bot.go
  echo "Ran $i times"
done

echo "Done"

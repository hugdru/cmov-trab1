PRAGMA foreign_keys = ON;
PRAGMA encoding = "UTF-8";

INSERT INTO voucher(idVoucher, idVoucherTemplate, idCustomer, idSale, code, gotVoucher, usedVoucher) VALUES
  (NULL,1,1,NULL,"111111111111111111111111111111111111111111111","111-111",NULL),
  (NULL,2,1,NULL,"999999999999999999999999999999999999999999999","999-999",NULL),
  (NULL,3,1,NULL,"555555555555555555555555555555555555555555555","555-555",NULL);
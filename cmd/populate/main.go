package main

import (
	"context"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/clients"
	"github.com/G1GACHADS/stashable-backend/config"
	"github.com/G1GACHADS/stashable-backend/core/logger"
	_ "github.com/joho/godotenv/autoload"
)

type address struct {
	province   string
	city       string
	streetName string
	zipCode    int
}

type room struct {
	name     string
	imageURL string
	width    float64
	height   float64
	length   float64
	price    float64
}

type warehouseData struct {
	name        string
	imageURL    string
	description string
	email       string
	phoneNumber string
	address     address
	categories  []int64
	rooms       []room
}

func main() {
	logger.Init(true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.New()
	clients, err := clients.New(ctx, config)
	if err != nil {
		logger.M.Fatal(err.Error())
	}

	b := backend.New(clients, config)

	b.RegisterUser(ctx, backend.User{
		FullName:    "John Doe",
		Email:       "user@mail.com",
		PhoneNumber: "0877824948548",
		Password:    "123123",
	}, backend.Address{
		Province:   "Jawa Barat",
		City:       "Bekasi",
		StreetName: "Jl. Belibis VI",
		ZipCode:    17421,
	})

	categoryChemical, _ := b.CategoryCreate(ctx, "Chemical")
	categoryElectricComponents, _ := b.CategoryCreate(ctx, "Electric")
	categoryFragileGlass, _ := b.CategoryCreate(ctx, "Fragile")
	categoryHeavyMaterials, _ := b.CategoryCreate(ctx, "Heavy Materials")

	var warehouseRealSampleData = []warehouseData{
		{
			name:        "Luntia Warehouse",
			imageURL:    "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662629437/stashable/photo-1565610222536-ef125c59da2e_x2wc6q.jpg",
			description: "Luntia Warehouse merupakan perusahaan yang bergerak di pergudangan dengan pengalaman selama 10 tahun. Berada di tempat yang strategis yaitu terletak di daerah Pluit, Jakarta Utara. Terdapat fasilitas yang memadai, seperti ruang tunggu yang dilengkapi ac, security/penjaga dan area bebas banjir",
			email:       "luntia@warehouse.com",
			phoneNumber: "089845674564",
			address: address{
				streetName: "Jalan Raden Saleh RT/RW: 01/05 No. 14",
				zipCode:    12345,
				city:       "Pluit",
				province:   "Jakarta Utara",
			},
			categories: []int64{categoryFragileGlass.ID, categoryElectricComponents.ID, categoryHeavyMaterials.ID},
			rooms: []room{
				{"Small Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662627185/stashable/istockphoto-1280808958-612x612_cutg5a.jpg", 1.2, 2.4, 3, 1000000},
				{"Medium Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662627291/stashable/205120-1600x1200-storageunit_qmf6ix.jpg", 1.5, 3, 4.5, 2000000},
			},
		},
		{
			name:        "Gudang Rune",
			imageURL:    "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662629385/stashable/photo-1633959592096-9d9a339b99fa_iw7cpz.jpg",
			description: "Tempat bebas banjir, tersedia ruang tunggu dengan luas area gudang 500m2. Perawatan setiap bulannya dan dijaga oleh security selama 24 jam/7 hari. Harga sewa ruangan mulai dari 1 juta sampai 3 juta saja. Cepat sebelum kehabisan",
			categories:  []int64{categoryChemical.ID, categoryFragileGlass.ID, categoryElectricComponents.ID, categoryHeavyMaterials.ID},
			email:       "rune@gudang.com",
			phoneNumber: "089845674564",
			address: address{
				streetName: "Jalan Damai 3 RT/RW: 01/10 No. 20",
				zipCode:    12345,
				city:       "Duri Kosambi",
				province:   "Jakarta Barat",
			},
			rooms: []room{
				{"Small Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662627185/stashable/istockphoto-1280808958-612x612_cutg5a.jpg", 1.2, 2.4, 3, 1500000},
				{"Large Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662632997/stashable/large_storage_room_wl3tz0.jpg", 3, 4, 6, 3000000},
			},
		},
		{
			name:        "IntearWarehouse",
			imageURL:    "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662629502/stashable/photo-1557761469-f29c6e201784_klw70x.jpg",
			description: "Melayani sewa tempat/ruangan dengan harga mulai dari 1 jutaan. Ruangan di lengkapi dengan faslitas sirkulasi udara yang baik, ruangan yang bersih, penjagaan 24 jam dan tempat bebas dari banjir.",
			categories:  []int64{categoryChemical.ID, categoryFragileGlass.ID, categoryElectricComponents.ID, categoryHeavyMaterials.ID},
			email:       "intear@warehouse.com",
			phoneNumber: "081212311561",
			address: address{
				streetName: "Jl. Kejawan Putih Mutiara No.17",
				zipCode:    12938,
				city:       "Surabaya",
				province:   "Jawa Timur",
			},
			rooms: []room{
				{"Small Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662627185/stashable/istockphoto-1280808958-612x612_cutg5a.jpg", 1.2, 2.4, 3, 1000000},
			},
		},
		{
			name:        "Gee Storage",
			imageURL:    "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662629502/stashable/photo-1557761469-f29c6e201784_klw70x.jpg",
			description: "Tersedia berbagai tempat untuk anda menyimpan barang di gudang kami. Setiap ruangan disediakan fasilitas penyaringan udara, kemanan pintu, tempat yang bersih dan tersedia rak jika dibutuhkan. Area sekitar yang bebas banjir serta mudah dijangkau dengan kendaraan besar maunpun kecil. ",
			categories:  []int64{categoryChemical.ID, categoryElectricComponents.ID, categoryHeavyMaterials.ID},
			email:       "gee@storage.com",
			phoneNumber: "084112346917",
			address: address{
				streetName: "Jl. Ahmad Yani No.88",
				zipCode:    581274,
				city:       "Surabaya",
				province:   "Jawa Timur",
			},
			rooms: []room{
				{"Large Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662632997/stashable/large_storage_room_wl3tz0.jpg", 3, 4, 6, 3500000},
			},
		},
		{
			name:        "Belly Rent Storage and Warehouse",
			imageURL:    "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662629502/stashable/photo-1557761469-f29c6e201784_klw70x.jpg",
			description: "Tersedia berbagai tempat untuk anda menyimpan barang di gudang kami. Setiap ruangan disediakan fasilitas penyaringan udara, kemanan pintu, tempat yang bersih dan tersedia rak jika dibutuhkan. Area sekitar yang bebas banjir serta mudah dijangkau dengan kendaraan besar maunpun kecil. ",
			categories:  []int64{categoryChemical.ID, categoryElectricComponents.ID, categoryHeavyMaterials.ID},
			email:       "gee@storage.com",
			phoneNumber: "084112346917",
			address: address{
				streetName: "Jl. Raya Darmo Permai Selatan No. 6-14",
				zipCode:    412475,
				city:       "Surabaya",
				province:   "Jawa Timur",
			},
			rooms: []room{
				{"Small Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662627185/stashable/istockphoto-1280808958-612x612_cutg5a.jpg", 1.2, 2.4, 3, 1200000},
				{"Medium Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662627291/stashable/205120-1600x1200-storageunit_qmf6ix.jpg", 1.5, 3, 4.5, 2200000},
				{"Large Room", "https://res.cloudinary.com/swcq1ct1pve4adyd2m4n/image/upload/v1662632997/stashable/large_storage_room_wl3tz0.jpg", 3, 4, 6, 3200000},
			},
		},
	}

	for i, warehouse := range warehouseRealSampleData {
		rooms := make([]backend.Room, len(warehouse.rooms))

		for idx, room := range warehouse.rooms {
			rooms[idx] = backend.Room{
				ImageURL: room.imageURL,
				Name:     room.name,
				Width:    room.width,
				Height:   room.height,
				Length:   room.length,
				Price:    room.price,
			}
		}

		logger.M.Debug(rooms[0].Price)
		err := b.WarehouseCreate(ctx, backend.WarehouseCreateInput{
			Warehouse: backend.Warehouse{
				Name:        warehouse.name,
				ImageURL:    warehouse.imageURL,
				Description: warehouse.description,
				BasePrice:   rooms[0].Price,
				Email:       warehouse.email,
				PhoneNumber: warehouse.phoneNumber,
			},
			Address: backend.Address{
				Province:   warehouse.address.province,
				City:       warehouse.address.city,
				StreetName: warehouse.address.streetName,
				ZipCode:    warehouse.address.zipCode,
			},
			Rooms:       rooms,
			CategoryIDs: warehouse.categories,
		})
		if err != nil {
			logger.M.Warnf("#%d: failed inserting\nreason:%v", i, err)
		}
	}

	logger.M.Info("Database populate process finished!")
}

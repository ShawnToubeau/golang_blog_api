package controllertests

//func TestCreateUser(t *testing.T) {
//	err := refreshUserTable()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	samples := []struct{
//		inputJSON string
//		statusCode int
//		nickname string
//		email string
//		errorMessage string
//	}{
//		{
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			201,
//			"Shawn",
//			"shawn@aol.com",
//			"",
//		},
//		{
//			`{"nickname": "Aria", "email": "shawn@aol.com", "password": "321"}`,
//			500,
//			"Aria",
//			"shawn@aol.com",
//			"Email Already Taken",
//		},
//		{
//			`{"nickname": "Shawn", "email": "aria@aol.com", "password": "321"}`,
//			500,
//			"Shawn",
//			"aria@aol.com",
//			"Nickname Already Taken",
//		},
//		{
//			`{"nickname": "Aria", "email": "", "password": "321"}`,
//			422,
//			"Aria",
//			"",
//			"Required Email",
//		},
//		{
//			`{"nickname": "", "email": "aria@aol.com", "password": "321"}`,
//			422,
//			"",
//			"aria@aol.com",
//			"Required Nickname",
//		},
//		{
//			`{"nickname": "Aria", "email": "", "password": "321"}`,
//			422,
//			"Aria",
//			"",
//			"Required Email",
//		},
//		{
//			`{"nickname": "Aria", "email": "aria@aol.com", "password": ""}`,
//			422,
//			"Aria",
//			"aria@aol.com",
//			"Required Password",
//		},
//	}
//
//	for _, v := range samples {
//		req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(v.inputJSON))
//		if err != nil {
//			t.Errorf("Error creating request: %v\n", err)
//		}
//		// Records server responses
//		rr := httptest.NewRecorder()
//		// Request handler
//		handler := http.HandlerFunc(server.InsertUser)
//		// Serve request
//		handler.ServeHTTP(rr, req)
//
//		responseMap := make(map[string]interface{})
//		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//		if err != nil {
//			fmt.Printf("Cannot convert to JSON: %v\n", err)
//		}
//		assert.Equal(t, rr.Code, v.statusCode)
//		if v.statusCode == 201 {
//			assert.Equal(t, responseMap["nickname"], v.nickname)
//			assert.Equal(t, responseMap["email"], v.email)
//		}
//		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
//			assert.Equal(t, responseMap["error"], v.errorMessage)
//		}
//	}
//}
//
//func TestGetUsers(t *testing.T) {
//	// Refresh table
//	err := refreshUserTable()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Seed table
//	_, err = seedUsers()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Create request
//	req, err := http.NewRequest("GET", "/users", nil)
//	if err != nil {
//		t.Errorf("Failed to form request: %v\n", err)
//	}
//	// Create request recorder and serve
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(server.FetchAllUsers)
//	handler.ServeHTTP(rr, req)
//	// Create user array and process response
//	var users []models.User
//	err = json.Unmarshal([]byte(rr.Body.String()), &users)
//	if err != nil {
//		log.Fatalf("Cannot convert to JSON: %v\n", err)
//	}
//
//	assert.Equal(t, rr.Code, http.StatusOK)
//	assert.Equal(t, len(users), 2)
//}
//
//func TestGetUserById(t *testing.T) {
//	// Refresh table
//	err := refreshUserTable()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Seed v
//	user, err := seedOneUser()
//	if err != nil {
//		log.Fatal(err)
//	}
//	userSample := []struct {
//		id string
//		statusCode int
//		nickname string
//		email string
//		password string
//	}{
//		{
//
//			strconv.Itoa(int(user.ID)),
//			200,
//			user.Nickname,
//			user.Email,
//			"",
//		},
//		{
//			"unknown",
//			400,
//			"",
//			"",
//			"",
//		},
//	}
//
//	for _, v := range userSample {
//		// Create request
//		req, err := http.NewRequest("GET", "/users", nil)
//		if err != nil {
//			t.Errorf("Failed to create request: %v\n", err)
//		}
//		// Set request params and test request
//		req = mux.SetURLVars(req, map[string]string{"id": v.id})
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.GetUserByID)
//		handler.ServeHTTP(rr, req)
//		// Process response
//		responseMap := make(map[string]interface{})
//		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//		if err != nil {
//			log.Fatalf("Cannot convert to JSON: %v\n", err)
//		}
//
//		assert.Equal(t, rr.Code, v.statusCode)
//
//		if v.statusCode == 200 {
//			assert.Equal(t, user.Nickname, responseMap["nickname"])
//			assert.Equal(t, user.Email, responseMap["email"])
//		}
//	}
//}
//
//func TestUpdateUser(t *testing.T) {
//	var AuthEmail, AuthPassword string
//	var AuthID uint32
//
//	// Refresh user table
//	err := refreshUserTable()
//	if err != nil {
//		log.Fatalf("Failed to refresh user table: %v\n", err)
//	}
//	// Seed one user
//	users, err := seedUsers()
//	if err != nil {
//		log.Fatalf("Failed to seed user table: %v\n", err)
//	}
//	// Retrieve first user
//	AuthID = users[0].ID
//	AuthEmail = users[0].Email
//	AuthPassword = "123"
//	// Login user to retrieve auth token
//	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
//	if err != nil {
//		log.Fatalf("Failed to login user: %v\n", err)
//	}
//	// Format token
//	tokenString := fmt.Sprintf("Bearer %v", token)
//
//	samples := []struct {
//		id string
//		updateJSON string
//		statusCode int
//		updateNickname string
//		updateEmail string
//		tokenGiven string
//		errorMessage string
//	} {
//		// OK
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			200,
//			"Shawn",
//			"shawn@aol.com",
//			tokenString,
//			"",
//		},
//		// Empty password field
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": ""}`,
//			422,
//			"Shawn",
//			"shawn@aol.com",
//			tokenString,
//			"Required Password",
//		},
//		// No auth token
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			401,
//			"Shawn",
//			"swan@aol.com",
//			"",
//			"Unauthorized",
//		},
//		// Wrong auth token
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			401,
//			"Shawn",
//			"swan@aol.com",
//			"wrong token",
//			"Unauthorized",
//		},
//		// Email taken
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			500,
//			"Shawn",
//			"aria@aol.com",
//			tokenString,
//			"Email Already Taken",
//		},
//		// Nickname taken
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			401,
//			"Aria",
//			"swan@aol.com",
//			tokenString,
//			"Nickname Already Taken",
//		},
//		// Email invalid
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			422,
//			"Shawn",
//			"",
//			tokenString,
//			"Invalid Email",
//		},
//		// Email field empty
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "Shawn", "email": "", "password": "123"}`,
//			422,
//			"",
//			"",
//			tokenString,
//			"Required Email",
//		},
//		// Nickname field empty
//		{
//			strconv.Itoa(int(AuthID)),
//			`{"nickname": "", "email": "shawn@aol.com", "password": "123"}`,
//			422,
//			"",
//			"",
//			tokenString,
//			"Required Nickname",
//		},
//		// No ID
//		{
//			"",
//			"",
//			400,
//			"",
//			"",
//			tokenString,
//			"",
//		},
//		// Using other user's token
//		{
//			strconv.Itoa(int(2)),
//			`{"nickname": "Shawn", "email": "shawn@aol.com", "password": "123"}`,
//			401,
//			"",
//			"",
//			tokenString,
//			"Unauthorized",
//		},
//	}
//
//	for _, v := range samples {
//		// Create request
//		urlStr := fmt.Sprintf("/users?token=%v", token)
//		req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(v.updateJSON))
//		if err != nil {
//			t.Errorf("Failed to create request: %v\n", err)
//		}
//		req = mux.SetURLVars(req, map[string]string{"id": v.id})
//
//
//
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.UpdateUserById)
//		req.Header.Set("Authorization", v.tokenGiven)
//
//		fmt.Printf("req %v\n", req)
//
//		handler.ServeHTTP(rr, req)
//
//		responseMap := make(map[string]interface{})
//		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//		if err != nil {
//			t.Errorf("Cannot convert to JSON: %v\n", err)
//		}
//
//		fmt.Printf("res map %v\n", responseMap)
//
//		assert.Equal(t, rr.Code, v.statusCode)
//
//		if v.statusCode == 200 {
//			assert.Equal(t, responseMap["nickname"], v.updateNickname)
//			assert.Equal(t, responseMap["email"], v.updateEmail)
//		}
//		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
//			assert.Equal(t, responseMap["error"], v.errorMessage)
//		}
//	}
//}
//
//func TestDeleteUser(t *testing.T) {
//	var AuthEmail, AuthPassword string
//	var AuthId uint32
//
//	// Refresh user table
//	err := refreshUserTable()
//	if err != nil {
//		log.Fatalf("Failed to refresh user table: %v\n", err)
//	}
//	// Seed users
//	users, err := seedUsers()
//	if err != nil {
//		log.Fatalf("Failed to see users: %v\n", err)
//	}
//	// Get first users credentials
//		AuthId = users[0].ID
//		AuthEmail = users[0].Email
//		AuthPassword = users[0].Password
//	// Login in the user to get their auth token
//	token, err := server.AuthenticateCredentials(AuthEmail, AuthPassword)
//	if err != nil {
//		log.Fatalf("Failed to login user: %v\n", err)
//	}
//	tokenString := fmt.Sprintf("Bearer %v", token)
//
//	userSample := []struct{
//		id string
//		tokenGiven string
//		stateCode int
//		errorMessage string
//	} {
//		{
//			strconv.Itoa(int(AuthId)),
//			tokenString,
//			204,
//			"",
//		},
//		{
//			strconv.Itoa(int(AuthId)),
//			"",
//			401,
//			"Unauthorized",
//		},
//		{
//			strconv.Itoa(int(AuthId)),
//			"Incorrect token",
//			401,
//			"Unauthorized",
//		},
//		{
//			"",
//			tokenString,
//			400,
//			"",
//		},
//		{
//			strconv.Itoa(int(2)),
//			tokenString,
//			401,
//			"Unauthorized",
//		},
//	}
//
//	for _, v := range userSample {
//		req, err := http.NewRequest("GET", "/users", nil)
//		if err != nil {
//			t.Errorf("Failed to create request: %v\n", err)
//		}
//		req = mux.SetURLVars(req, map[string]string{"id": v.id})
//		rr := httptest.NewRecorder()
//		handler := http.HandlerFunc(server.DeleteUserById)
//
//		req.Header.Set("Authorization", v.tokenGiven)
//
//		handler.ServeHTTP(rr, req)
//		assert.Equal(t, rr.Code, v.stateCode)
//
//		if v.stateCode == 401 && v.errorMessage != "" {
//			responseMap := make(map[string]interface{})
//			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
//			if err != nil {
//				t.Errorf("Cannot convert to json: %v\n", err)
//			}
//			assert.Equal(t, responseMap["error"], v.errorMessage)
//		}
//	}
//}

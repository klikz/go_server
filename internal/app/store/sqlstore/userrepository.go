package sqlstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"premier_api/internal/app/model"
	"premier_api/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

func setPin(param, addres string) (interface{}, error) {
	response, err := http.PostForm(addres, url.Values{
		"status": {param}})
	if err != nil {
		return nil, err

	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("SetPin err: ", err)
		return nil, err
	}
	return string(body), nil
}

func (r *UserRepository) debitFromLine(modelId, lineId int) (interface{}, error) {
	type Debit struct {
		Component_id int
		Quantity     float64
	}
	fmt.Println("debit modelId: ", modelId, "debit lineId: ", lineId)

	rows, err := r.store.db.Query(fmt.Sprintf("select t.component_id, t.quantity  from models.\"%d\" t, public.components c where t.component_id = c.id and c.\"checkpoint\" = %d", modelId, lineId))
	if err != nil {
		fmt.Println("debitFromLinel err: ", err)
		return nil, err
	}
	defer rows.Close()

	var debits []Debit
	// fmt.Println(rows)

	for rows.Next() {
		var debit Debit
		if err := rows.Scan(&debit.Component_id, &debit.Quantity); err != nil {
			fmt.Println("debitFromLine2 err: ", err)
			return debits, err
		}
		debits = append(debits, debit)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("debitFromLine3 err: ", err)
		return debits, err
	}
	for _, x := range debits {
		_, err := r.store.db.Exec(fmt.Sprintf("update checkpoints.\"%d\" set quantity = quantity - %f where component_id = %d", lineId, x.Quantity, x.Component_id))
		if err != nil {
			fmt.Println("error in debit: ", err)
			return nil, err
		}

	}
	return nil, nil
}

func CheckLaboratory(serial string) (string, error) {
	response, err := http.PostForm("http://192.168.5.250:3002/labinfo", url.Values{
		"serial": {serial}})
	if err != nil {
		fmt.Println("CheckLaboratory err: ", err)
		return "", err

	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("CheckLaboratory2 err: ", err)
		return "", err
	}

	fmt.Println("Laboratory: ", string(body))
	return string(body), nil
}
func (r *UserRepository) GetLines() (interface{}, error) {

	type Lines struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	rows, err := r.store.db.Query("select c.id, c.\"name\"  from checkpoints c ")
	if err != nil {
		fmt.Println("GetLines err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []Lines

	for rows.Next() {
		var comp Lines
		if err := rows.Scan(&comp.ID,
			&comp.Name); err != nil {
			fmt.Println("GetLines2 err: ", err)
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetLines3 err: ", err)
		return last, err
	}

	return last, nil
}
func (r *UserRepository) DeleteDefectsTypes(id int) (interface{}, error) {
	rows, err := r.store.db.Query("delete from defects where id = $1", id)
	if err != nil {
		fmt.Println("deletedefectsTypes err: ", err)
		return nil, err
	}
	defer rows.Close()
	return nil, nil
}

func (r *UserRepository) AddDefectsTypes(id int, name string) (interface{}, error) {
	rows, err := r.store.db.Query("insert into defects (defect_name, line_id) values ($1, $2)", name, id)
	if err != nil {
		fmt.Println("AdddefectsTypes err: ", err)
		return nil, err
	}
	defer rows.Close()
	return nil, nil
}
func (r *UserRepository) GetDefectsTypes() (interface{}, error) {

	type defectsTypes struct {
		ID          int    `json:"id"`
		Defect_name string `json:"defect_name"`
		Line_id     string `json:"line_id"`
		Name        string `json:"name"`
	}

	rows, err := r.store.db.Query("select r.id, r.defect_name, r.line_id, c.\"name\"  from defects r, checkpoints c where c.id = r.line_id order by line_id")
	if err != nil {
		fmt.Println("GetdefectsTypes err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []defectsTypes

	for rows.Next() {
		var comp defectsTypes
		if err := rows.Scan(&comp.ID,
			&comp.Defect_name,
			&comp.Line_id,
			&comp.Name); err != nil {
			fmt.Println("GetdefectsTypes2 err: ", err)
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetdefectsTypes3 err: ", err)
		return last, err
	}

	return last, nil
}
func (r *UserRepository) GetStatus(line int) (interface{}, error) {
	type Status struct {
		Status byte `json:"status"`
	}
	var last Status
	err := r.store.db.QueryRow("select c.status from checkpoints c where c.id = $1", line).Scan(&last.Status)
	if err != nil {
		fmt.Println("GetStatus err: ", err)
		return nil, err
	}

	return last, nil
}
func (r *UserRepository) GetPackingTodaySerial() (interface{}, error) {

	type PackingTodaySerial struct {
		Serial  string `json:"serial"`
		Packing string `json:"packing"`
		Time    string `json:"time"`
	}
	currentTime := time.Now()
	rows, err := r.store.db.Query("select serial, packing, to_char(\"time\" , 'DD-MM-YYYY HH24:MI') \"time\" from packing where \"time\"::date=to_date($1, 'YYYY-MM-DD') order by serial", currentTime)
	if err != nil {
		fmt.Println("GetPackingTodaySerial err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []PackingTodaySerial

	for rows.Next() {
		var comp PackingTodaySerial
		if err := rows.Scan(&comp.Serial,
			&comp.Packing,
			&comp.Time); err != nil {
			fmt.Println("GetPackingTodaySerial2 err: ", err)
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetPackingTodaySerial3 err: ", err)
		return last, err
	}

	return last, nil
}

func (r *UserRepository) PackingSerialInput(serial, packing string) (interface{}, error) {

	type Laboratory struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Duration  string `json:"duration"`
		Model     string `json:"model"`
		Result    string `json:"result"`
	}

	res, err := CheckLaboratory(serial)
	if err != nil {
		fmt.Println("PackingSerialInput1 err: ", err)
		return nil, nil
	}

	s := string(res)
	data := Laboratory{}
	json.Unmarshal([]byte(s), &data)
	if data.Result == "No data" {
		return "", errors.New("laboratoriyada muammo")
	}
	type ModelId struct {
		id int
	}
	var modelId ModelId
	var serialSlice = serial[0:6]
	//check address of station
	//check model
	if err := r.store.db.QueryRow("select m.id from models m where m.code = $1", serialSlice).Scan(&modelId.id); err != nil {
		fmt.Println("Packing serial intput error: ")
		return nil, errors.New("serial xato")
	}

	rows, err := r.store.db.Query("insert into packing (serial, packing, model_id) values ($1, $2, $3)", serial, packing, modelId.id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type respond struct {
		Result string `json:"result"`
	}

	result := respond{
		Result: "ok",
	}

	return result, nil
}

func (r *UserRepository) GetPackingTodayModels() (interface{}, error) {

	type PackingTodayModels struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    int    `json:"count"`
	}
	currentTime := time.Now()
	rows, err := r.store.db.Query("select p.model_id, m.\"name\", COUNT(*) FROM packing p, models m where p.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and m.id = p.model_id group by m.\"name\", p.model_id", currentTime)
	if err != nil {
		fmt.Println("GetPackingTodayModels1 err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []PackingTodayModels

	for rows.Next() {
		var comp PackingTodayModels
		if err := rows.Scan(&comp.Model_id, &comp.Name, &comp.Count); err != nil {
			fmt.Println("GetPackingTodayModels2 err: ", err)
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetPackingTodayModels3 err: ", err)
		return nil, err
	}
	return last, nil
}
func (r *UserRepository) GetPackingToday() (interface{}, error) {

	type PackingToday struct {
		Count int `json:"count"`
	}
	currentTime := time.Now()
	var last PackingToday
	err := r.store.db.QueryRow("select count(*) from packing where \"time\"::date=to_date($1, 'YYYY-MM-DD')", currentTime).Scan(&last.Count)
	if err != nil {
		fmt.Println("GetPackingToday err: ", err)
		return nil, err
	}
	return last, nil
}
func (r *UserRepository) GetPackingLast() (interface{}, error) {

	type PackingLast struct {
		ID      int    `json:"id"`
		Serial  string `json:"serial"`
		Packing string `json:"packing"`
		Time    string `json:"time"`
	}

	rows, err := r.store.db.Query("select p.id, p.serial, p.packing, to_char(p.\"time\" , 'DD-MM-YYYY HH24:MI') \"time\" from packing p ORDER BY p.\"time\" DESC LIMIT 3")
	if err != nil {
		fmt.Println("GetPackingLast err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []PackingLast

	for rows.Next() {
		var comp PackingLast
		if err := rows.Scan(&comp.ID,
			&comp.Serial,
			&comp.Packing,
			&comp.Time); err != nil {
			fmt.Println("GetPackingLast2 err: ", err)
			return nil, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetPackingLast3 err: ", err)
		return last, err
	}
	return last, nil
}
func (r *UserRepository) SerialInput(line int, serial string) (interface{}, error) {
	type InputInfo struct {
		id      int
		address string
	}
	var modelInfo InputInfo
	var serialSlice = serial[0:6]
	//check address of station
	if err := r.store.db.QueryRow("select address from checkpoints where id = $1", line).Scan(&modelInfo.address); err != nil {
		return nil, errors.New("sector address topilmadi")
	}
	//check model
	if err := r.store.db.QueryRow("select m.id from models m where m.code = $1", serialSlice).Scan(&modelInfo.id); err != nil {
		req, err := setPin("0", modelInfo.address)
		if err != nil {
			fmt.Println("SerialInput Setpin err: ", err)
			return nil, err
		}
		fmt.Println("from raspberry: ", req)
		return nil, errors.New("serial xato")
	}
	type product_id struct {
		id int
	}
	var prod_id product_id
	//check stations before
	type CheckStation struct {
		product_id int
	}
	switch line {
	//check sborka for ppu
	case 10:
		check := &CheckStation{}
		if err := r.store.db.QueryRow("select product_id from production p where serial = $1 and  checkpoint_id = $2", serial, 2).Scan(&check.product_id); err != nil {
			req, err := setPin("0", modelInfo.address)
			if err != nil {
				fmt.Println("SerialInput Setpin err: ", err)
				return nil, err
			}
			fmt.Println("from raspberry: ", req)
			return nil, errors.New("Sborkada reg qilinmagan")
		}
		// defer result.Close()
		fmt.Println("check sborka in ppu result: ", check)
		//select  product_id  from production p where serial = 'A1FA010003674' and checkpoint_id = 2
	}

	// check production to serial
	if err := r.store.db.QueryRow("select product_id from production p where serial = $1 and  checkpoint_id = $2", serial, line).Scan(&prod_id.id); err == nil {
		if _, err := r.store.db.Exec("update production set updated = now() where product_id = $1", prod_id.id); err != nil {
			return nil, err
		}
		req, err := setPin("1", modelInfo.address)
		if err != nil {
			fmt.Println("SerialInput Setpin2 err: ", err)
			return nil, err
		}
		fmt.Println("from raspberry: ", req)
		return nil, errors.New("serial kiritilgan")
	} else {
		rows, err := r.store.db.Query("insert into production (model_id, serial, checkpoint_id) values ($1, $2, $3)", modelInfo.id, serial, line)

		if err != nil {
			fmt.Println("SerialInput3 Setpin err: ", err)
		}
		defer rows.Close()
		res, err2 := r.debitFromLine(modelInfo.id, line)
		if err2 != nil {
			fmt.Println("inputSerial debit err: ", err2)
		}
		req, err := setPin("1", modelInfo.address)
		if err != nil {
			fmt.Println("SerialInput rasp err: ", err)
			return nil, err
		}

		fmt.Println("from raspberry: ", req)
		fmt.Println("inputSerial debit res: ", res)

	}
	type respond struct {
		Result string `json:"result"`
	}

	result := respond{
		Result: "ok",
	}

	return result, nil
}

func (r *UserRepository) ComponentsAll() (interface{}, error) {
	rows, err := r.store.db.Query("select c.available, c.id, c.code, c.\"name\", c2.\"name\" as Checkpoint, c2.id as checkpoint_id,  c.unit, c.specs, c.photo, to_char(c.\"time\", 'DD-MM-YYYY HH24:MI') \"time\", t.\"name\" as type, t.id as type_id, c.weight from components c join checkpoints c2 on c2.id = c.\"checkpoint\" join \"types\" t on t.id  = c.\"type\" order by c.code")
	if err != nil {
		fmt.Println("ComponentsAll err: ", err)
		return nil, err
	}
	defer rows.Close()

	var components []model.Component
	// fmt.Println(rows)

	for rows.Next() {
		var comp model.Component
		if err := rows.Scan(&comp.Available, &comp.ID, &comp.Code,
			&comp.Name, &comp.Checkpoint, &comp.Checkpoint_id, &comp.Unit, &comp.Specs, &comp.Photo, &comp.Time, &comp.Type, &comp.Type_id, &comp.Weight); err != nil {
			fmt.Println("ComponentsAll2 err: ", err)
			return components, err
		}
		components = append(components, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("ComponentsAll3 err: ", err)
		return components, err
	}
	return components, nil
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		fmt.Println("Create err: ", err)
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		fmt.Println("Create BeforeCreate err: ", err)
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, role FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Find err: ", err)
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, role FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("FindByEmail err: ", err)
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetLast(line int) (interface{}, error) {

	type Last struct {
		Serial        string `json:"serial"`
		Model_id      int    `json:"model_id"`
		Model         string `json:"model"`
		Checkpoint_id int    `json:"checkpoint_id"`
		Line          string `json:"line"`
		Product_id    int    `json:"product_id"`
		Time          string `json:"time"`
	}

	rows, err := r.store.db.Query("select p.serial, p.model_id, m.\"name\" as model, p.checkpoint_id, c.\"name\" as line, p.product_id,  to_char(p.\"time\" , 'DD-MM-YYYY HH24:MI') \"time\" from production p, checkpoints c, models m where m.id = p.model_id and c.id = p.checkpoint_id and p.checkpoint_id = $1 ORDER BY p.\"time\" DESC LIMIT 2", line)
	if err != nil {
		fmt.Println("GetLast1 err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []Last

	for rows.Next() {
		var comp Last
		if err := rows.Scan(&comp.Serial,
			&comp.Model_id,
			&comp.Model,
			&comp.Checkpoint_id,
			&comp.Line,
			&comp.Product_id,
			&comp.Time); err != nil {
			fmt.Println("GetLast2 err: ", err)
			return last, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetLast3 err: ", err)
		return last, err
	}
	return last, nil

}

func (r *UserRepository) GetToday(line int) (interface{}, error) {

	type Count struct {
		Count int `json:"count"`
	}

	currentTime := time.Now()

	rows, err := r.store.db.Query("select count(*) from production where checkpoint_id = $1 and \"time\"::date=to_date($2, 'YYYY-MM-DD')", line, currentTime)
	if err != nil {
		fmt.Println("GetToday err: ", err)
		return nil, err
	}

	defer rows.Close()
	var last []Count

	count := Count{}

	for rows.Next() {
		var comp Count
		if err := rows.Scan(&count.Count); err != nil {
			fmt.Println("GetToday2 err: ", err)
			return last, err
		}
		last = append(last, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetToday3 err: ", err)
		return last, err
	}

	return count, nil
}
func (r *UserRepository) GetTodayModels(line int) (interface{}, error) {

	type ByModel struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    string `json:"count"`
	}

	currentTime := time.Now()

	rows, err := r.store.db.Query("select p.model_id, m.\"name\", COUNT(*) FROM production p, models m where p.checkpoint_id = $1 and p.\"time\"::date>=to_date($2, 'YYYY-MM-DD') and m.id = p.model_id group by m.\"name\", p.model_id", line, currentTime)
	if err != nil {
		fmt.Println("GetTodayModels err: ", err)
		return nil, err
	}

	defer rows.Close()
	var byModel []ByModel

	for rows.Next() {
		var comp ByModel
		if err := rows.Scan(&comp.Model_id,
			&comp.Name,
			&comp.Count); err != nil {
			fmt.Println("GetTodayModels2 err: ", err)
			return byModel, err
		}
		byModel = append(byModel, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetTodayModels3 err: ", err)
		return byModel, err
	}
	return byModel, nil
}
func (r *UserRepository) GetSectorBalance(line int) (interface{}, error) {

	type Balance struct {
		Component_id int     `json:"component_id"`
		Code         string  `json:"code"`
		Quantity     float32 `json:"quantity"`
		Name         string  `json:"name"`
	}

	rows, err := r.store.db.Query(fmt.Sprintf("select t.component_id, c.code,  t.quantity, c.\"name\" from checkpoints.\"%d\" t, components c where t.component_id = c.id ORDER BY t.quantity", line))
	if err != nil {
		fmt.Println("GetSectorBalance err: ", err)
		return nil, err
	}

	defer rows.Close()
	var balance []Balance

	for rows.Next() {
		var comp Balance
		if err := rows.Scan(&comp.Component_id,
			&comp.Code,
			&comp.Quantity,
			&comp.Name); err != nil {
			fmt.Println("GetSectorBalance2 err: ", err)
			return balance, err
		}
		balance = append(balance, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetSectorBalance3 err: ", err)
		return balance, err
	}
	return balance, nil
}

func (r *UserRepository) GetByDate(date1, date2 string, line int) (interface{}, error) {

	type Count struct {
		Count int `json:"count"`
	}

	count := Count{}
	switch line {
	case 13:
		rows, err := r.store.db.Query(fmt.Sprintf(`select count(*) from packing where "time"::date>=to_date('%s', 'YYYY-MM-DD') and "time"::date<=to_date('%s', 'YYYY-MM-DD')`, date1, date2))
		if err != nil {
			fmt.Println("GetToday err: ", err)
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&count.Count); err != nil {
				fmt.Println("GetToday2 err: ", err)
				return count, err
			}

		}
		if err = rows.Err(); err != nil {
			fmt.Println("GetToday3 err: ", err)
			return count, err
		}
		break
	default:
		rows, err := r.store.db.Query("select count(*) from production where \"time\"::date>=to_date($1, 'YYYY-MM-DD') and \"time\"::date<=to_date($2, 'YYYY-MM-DD') and checkpoint_id = $3", date1, date2, line)
		if err != nil {
			fmt.Println("GetToday err: ", err)
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&count.Count); err != nil {
				fmt.Println("GetToday2 err: ", err)
				return count, err
			}
		}
		if err = rows.Err(); err != nil {
			fmt.Println("GetToday3 err: ", err)
			return count, err
		}
		break

	}

	return count, nil
}

func (r *UserRepository) GetByDateModels(date1, date2 string, line int) (interface{}, error) {

	type ByModel struct {
		Model_id int    `json:"model_id"`
		Name     string `json:"name"`
		Count    string `json:"count"`
	}
	var byModel []ByModel

	switch line {
	case 13:
		rows, err := r.store.db.Query("select p.model_id, m.\"name\", COUNT(*) FROM packing p, models m where p.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p.model_id group by m.\"name\", p.model_id", date1, date2)
		if err != nil {
			fmt.Println("GetTodayModels err: ", err)
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var comp ByModel
			if err := rows.Scan(&comp.Model_id,
				&comp.Name,
				&comp.Count); err != nil {
				fmt.Println("GetTodayModels2 err: ", err)
				return byModel, err
			}
			byModel = append(byModel, comp)
		}
		if err = rows.Err(); err != nil {
			fmt.Println("GetTodayModels3 err: ", err)
			return byModel, err
		}
		break
	default:
		rows, err := r.store.db.Query("select p.model_id, m.\"name\", COUNT(*) FROM production p, models m where p.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and checkpoint_id = $3 and m.id = p.model_id group by m.\"name\", p.model_id", date1, date2, line)
		if err != nil {
			fmt.Println("GetTodayModels err: ", err)
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var comp ByModel
			if err := rows.Scan(&comp.Model_id,
				&comp.Name,
				&comp.Count); err != nil {
				fmt.Println("GetTodayModels2 err: ", err)
				return byModel, err
			}
			byModel = append(byModel, comp)
		}
		if err = rows.Err(); err != nil {
			fmt.Println("GetTodayModels3 err: ", err)
			return byModel, err
		}
		break
	}

	return byModel, nil
}
func (r *UserRepository) GetByDateSerial(date1, date2 string) (interface{}, error) {
	type Serial struct {
		Serial string `json:"serial"`
		Model  string `json:"model"`
		Time   string `json:"time"`
		Sector string `json:"sector"`
	}
	var serial []Serial
	// rows, err := r.store.db.Query("(select p.serial, m.\"name\" as model, p.\"time\", c.\"name\" as sector  from packing p, models m, checkpoints c  where p.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p.model_id and c.id = p.checkpoint_id  order by p.model_id) union ALL (select p2.serial, m.\"name\" as model, p2.\"time\", c.\"name\" as sector  from production p2, models m, checkpoints c where p2.\"time\"::date>=to_date($1, 'YYYY-MM-DD') and p2.\"time\"::date<=to_date($2, 'YYYY-MM-DD') and m.id = p2.model_id and c.id = p2.checkpoint_id order by p2.model_id, p2.checkpoint_id)", date1, date2)
	rows, err := r.store.db.Query(fmt.Sprintf(`
	(select p.serial, m."name" as model, to_char(p."time" , 'DD-MM-YYYY HH24:MI') "time" , c."name" as sector  from packing p, models m, checkpoints c  
	where p."time"::date>=to_date('%s', 'YYYY-MM-DD') and p."time"::date<=to_date('%s', 'YYYY-MM-DD') and m.id = p.model_id and c.id = p.checkpoint_id  order by p.model_id) 
	union ALL 
	(select p2.serial, m."name" as model, to_char(p2."time" , 'DD-MM-YYYY HH24:MI') "time", c."name" as sector  from production p2, models m, checkpoints c 
	where p2."time"::date>=to_date('%s', 'YYYY-MM-DD') and p2."time"::date<=to_date('%s', 'YYYY-MM-DD') and m.id = p2.model_id and c.id = p2.checkpoint_id order by p2.model_id, p2.checkpoint_id)`,
		date1, date2, date1, date2))
	if err != nil {
		fmt.Println("GetToday err: ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comp Serial
		if err := rows.Scan(&comp.Serial, &comp.Model, &comp.Time, &comp.Sector); err != nil {
			fmt.Println("GetToday2 err: ", err)
			return serial, err
		}
		serial = append(serial, comp)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetToday3 err: ", err)
		return serial, err
	}
	return serial, nil
}

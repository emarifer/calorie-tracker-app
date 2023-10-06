import { useEffect, useState } from "react"
import {
    Button,
    Container,
    FormGroup,
    FormLabel,
    FormControl,
    Modal
} from "react-bootstrap";
import { Entry } from "./single-entry.component";

import 'bootstrap/dist/css/bootstrap.css'

export const Entries = () => {
    const [entries, setEntries] = useState([])
    const [addNewEntry, setAddNewEntry] = useState(false)
    const [newEntry, setNewEntry] = useState({
        dish: '',
        ingredients: '',
        calories: '',
        fat: 0,
    });
    const [refreshData, setRefreshData] = useState(false)
    const [changeIngredients, setChangeIngredients] = useState({
        change: false,
        id: '',
        ingredientsList: '',
    })
    const [newIngredientName, setNewIngredientName] = useState("")
    const [changeEntry, setChangeEntry] = useState({
        change: false,
        entry: {},
    })

    const getAllEntries = async () => {
        const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        // console.log(timeZone);
        const response = await fetch('/api/entries', {
            headers: {
                'X-TimeZone': timeZone,
            }
        });
        const data = await response.json();

        // console.log(data)
        setEntries(data)
    }

    const addSingleEntry = async () => {
        setAddNewEntry(false)

        const response = await fetch('/api/entry/create', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "dish": newEntry.dish,
                "ingredients": newEntry.ingredients,
                "calories": newEntry.calories,
                "fat": parseFloat(newEntry.fat),
            }),
        });

        // console.log(response.status)
        if (response.status == 201) setRefreshData(true)
    }

    const deleteEntry = async (id) => {
        const result = confirm('Are you sure to eliminate this dish?');
        if (!result) return;

        const response = await fetch(`/api/entry/${id}`, { method: "DELETE" });

        if (response.status == 200) setRefreshData(true)
    }

    const changeIngredientsForEntry = async () => {
        changeIngredients.change = false;
        const url = `/api/ingredient/update/${changeIngredients.id}`;
        // agrega ingredientes a la lista
        const ingredients = `${changeIngredients.ingredientsList}, ${newIngredientName}`;

        const response = await fetch(url, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "ingredients": ingredients, // newIngredientName
            }),
        });

        // console.log(response.status)
        if (response.status == 200) setRefreshData(true)
    }

    const emptyFormCheck = (dish, ingredients, calories, fat) => {
        return (
            dish === '' && ingredients === '' && calories === '' && fat === 0
        )
    }

    const changeSingleEntry = async () => {
        const url = `/api/entry/update/${changeEntry.entry.ID}`;
        // se verifica el formulario no esté completamente vacío
        // porque el backend enviaría un "404" dado que mongodb
        // no actualiza algo con los mismo datos que ya tiene
        if (!emptyFormCheck(
            newEntry.dish.trim(),
            newEntry.ingredients.trim(),
            newEntry.calories.trim(),
            newEntry.fat,
        )) {
            // se verifica que los campos no estén vacíos uno a uno;
            // si alguno lo está, se envía el valor anterior que había
            const dish = newEntry.dish.trim() === ''
                ? changeEntry.entry.dish // valor anterior
                : newEntry.dish.trim(); // valor del campo del formulario
            const ingredients = newEntry.ingredients.trim() === ''
                ? changeEntry.entry.ingredients
                : newEntry.ingredients.trim();
            const calories = newEntry.calories.trim() === ''
                ? changeEntry.entry.calories
                : newEntry.calories.trim();
            const fat = newEntry.fat === 0
                ? parseFloat(changeEntry.entry.fat)
                : parseFloat(newEntry.fat);

            // finalmente se verifica que NINGUNO de los valores
            // que se van a enviar sean iguales a los anteriores
            // por las mismas razones dadas más arriba.
            // Si el backend manda un 404 es que el "entry" no está en la DB
            if (!(dish === changeEntry.entry.dish &&
                ingredients === changeEntry.entry.ingredients &&
                calories === changeEntry.entry.calories &&
                parseFloat(fat) === parseFloat(changeEntry.entry.fat))) {
                const response = await fetch(url, {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        "dish": dish,
                        "ingredients": ingredients,
                        "calories": calories,
                        "fat": fat,
                    }),
                });

                // console.log(response.status)
                if (response.status == 200) setRefreshData(true)
            }

            // se resetea el estado del formulario
            setNewEntry({
                dish: '',
                ingredients: '',
                calories: '',
                fat: 0,
            });
        }
    }

    useEffect(() => {
        getAllEntries()
    }, [])

    if (refreshData) {
        setRefreshData(false)
        getAllEntries()
    }

    return (
        <>
            <Container>
                <Button onClick={() => setAddNewEntry(true)}>
                    Track today's calories
                </Button>
            </Container>
            <Container>
                {entries && entries.map((entry) => (
                    <Entry
                        key={entry.ID}
                        entryData={entry}
                        deleteEntry={deleteEntry}
                        setChangeIngredients={setChangeIngredients}
                        setChangeEntry={setChangeEntry}
                    />
                ))}
            </Container>

            {/* Add Entry Modal */}
            <Modal
                show={addNewEntry}
                onHide={() => setAddNewEntry(false)}
                centered>
                <Modal.Header closeButton>
                    <Modal.Title>Add Calorie Entry</Modal.Title>
                </Modal.Header>

                <Modal.Body>
                    <FormGroup>
                        <FormLabel>dish</FormLabel>
                        <FormControl onChange={(e) => newEntry.dish = e.target.value} required autoFocus>
                        </FormControl>

                        <FormLabel>ingredients</FormLabel>
                        <FormControl onChange={(e) => newEntry.ingredients = e.target.value} required>
                        </FormControl>

                        <FormLabel>calories</FormLabel>
                        <FormControl onChange={(e) => newEntry.calories = e.target.value} required>
                        </FormControl>

                        <FormLabel>fat</FormLabel>
                        <FormControl type="number" onChange={(e) => newEntry.fat = e.target.value} required>
                        </FormControl>
                    </FormGroup>

                    <Button onClick={addSingleEntry}>
                        Add
                    </Button>
                    <Button onClick={() => setAddNewEntry(false)}>
                        Cancel
                    </Button>
                </Modal.Body>
            </Modal>

            {/* Change Ingredients Modal */}
            <Modal
                show={changeIngredients.change}
                onHide={() => setChangeIngredients({
                    change: false,
                    id: '',
                    ingredientsList: '',
                })}
                centered>
                <Modal.Header closeButton>
                    <Modal.Title>Change Ingredients</Modal.Title>
                </Modal.Header>

                <Modal.Body>
                    <FormGroup>
                        <FormLabel>new ingredients</FormLabel>
                        <FormControl onChange={
                            (e) => setNewIngredientName(e.target.value)
                        } required autoFocus>
                        </FormControl>
                    </FormGroup>

                    <Button onClick={changeIngredientsForEntry}>
                        Change
                    </Button>
                    <Button onClick={() => setChangeIngredients({
                        change: false,
                        id: '',
                        ingredientsList: '',
                    })}>
                        Cancel
                    </Button>
                </Modal.Body>
            </Modal>

            {/* Change Entry Modal */}
            <Modal
                show={changeEntry.change}
                onHide={() => setChangeEntry({
                    change: false,
                    entry: {},
                })}
                centered>
                <Modal.Header closeButton>
                    <Modal.Title>Change Entry</Modal.Title>
                </Modal.Header>

                <Modal.Body>
                    <FormGroup>
                        <FormLabel>dish</FormLabel>
                        <FormControl onChange={(e) => newEntry.dish = e.target.value} autoFocus>
                        </FormControl>

                        <FormLabel>ingredients</FormLabel>
                        <FormControl onChange={(e) => newEntry.ingredients = e.target.value}>
                        </FormControl>

                        <FormLabel>calories</FormLabel>
                        <FormControl onChange={(e) => newEntry.calories = e.target.value}>
                        </FormControl>

                        <FormLabel>fat</FormLabel>
                        <FormControl type="number" onChange={(e) => newEntry.fat = e.target.value}>
                        </FormControl>
                    </FormGroup>

                    <Button onClick={() => {
                        changeSingleEntry();
                        setChangeEntry({
                            change: false,
                            entry: {},
                        });
                    }}>
                        Change Entry
                    </Button>
                    <Button onClick={() => setChangeEntry({
                        change: false,
                        entry: {},
                    })}>
                        Cancel
                    </Button>
                </Modal.Body>
            </Modal>
        </>
    )
}

/* Styles inline in React. VER:
https://www.pluralsight.com/guides/inline-styling-with-react
*/
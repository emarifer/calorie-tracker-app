import dateFormat from "dateformat"
import { Button, Card, Col, Row } from "react-bootstrap"

export const Entry = ({
    entryData,
    deleteEntry,
    setChangeIngredients,
    setChangeEntry,
}) => {
    const convertDateTime = (datetime) => {
        return dateFormat(datetime, 'dd mmm yyyy, HH:MM')
    }

    const changeIngredients = () => {
        setChangeIngredients({
            change: true,
            id: entryData.ID,
            ingredientsList: entryData.ingredients,
        });
    }

    const changeEntry = () => {
        setChangeEntry({
            change: true,
            entry: entryData,
        });
    }

    return (
        <>
            <Card>
                <Row>
                    <Col>
                        Dish: {entryData && entryData.dish}
                        <p style={{
                            fontSize: '10px', width: '100px', fontWeight: 'bold', marginTop: '8px'
                        }}>
                            ({entryData && convertDateTime(entryData.createdAt)})
                        </p>
                    </Col>
                    <Col>Ingredients: {entryData && entryData.ingredients}</Col>
                    <Col>Calories: {entryData && entryData.calories}</Col>
                    <Col>Fat: {entryData && entryData.fat}</Col>
                    <Col>
                        <Button onClick={() => deleteEntry(entryData.ID)}>
                            Delete Entry
                        </Button>
                    </Col>
                    <Col>
                        <Button onClick={changeIngredients}>
                            Change Ingredients
                        </Button>
                    </Col>
                    <Col>
                        <Button onClick={changeEntry}>
                            Change Entry
                        </Button>
                    </Col>
                </Row>
            </Card>
        </>
    );
}

/* Formatear el datetime como RFC822Z o similar
https://stackoverflow.com/questions/64590320/generate-a-rfc2822-datetime-string-node-js
https://www.npmjs.com/package/dateformat
*/
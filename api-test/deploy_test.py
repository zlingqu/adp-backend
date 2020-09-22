import unittest
import requests

space_addr = 'http://localhost:80/api/v1/'


class ADPAPITestCase(unittest.TestCase):
    def test_put_space(self):
        resp = requests.put(space_addr + "space/112", json={"name": "aixh", "owner": "liaolonglong"})
        print(resp.text)


if __name__ == '__main__':
    unittest.main()

import unittest

import requests

addr = 'http://localhost:80/api/v1/'


class ADPAPITestCase(unittest.TestCase):
    def test_put_space(self):
        resp = requests.put(addr + "space/112", json={"name": "aixh", "owner": "liaolonglong"})
        print(resp.text)

    def test_post_result(self):
        for i in range(10):
            resp = requests.post(addr + "result", json={"name": "cp-wechat", "deploy_env": "dev", "version": "last"})
            print(resp.text)


if __name__ == '__main__':
    unittest.main()

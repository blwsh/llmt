<?php

// autoload_real.php @generated by Composer

class ComposerAutoloaderInitcfb7004a461833b2c554ede0980f1246
{
    private static $loader;

    public static function loadClassLoader($class)
    {
        if ('Composer\Autoload\ClassLoader' === $class) {
            require __DIR__ . '/ClassLoader.php';
        }
    }

    /**
     * @return \Composer\Autoload\ClassLoader
     */
    public static function getLoader()
    {
        if (null !== self::$loader) {
            return self::$loader;
        }

        spl_autoload_register(array('ComposerAutoloaderInitcfb7004a461833b2c554ede0980f1246', 'loadClassLoader'), true, true);
        self::$loader = $loader = new \Composer\Autoload\ClassLoader(\dirname(__DIR__));
        spl_autoload_unregister(array('ComposerAutoloaderInitcfb7004a461833b2c554ede0980f1246', 'loadClassLoader'));

        require __DIR__ . '/autoload_static.php';
        call_user_func(\Composer\Autoload\ComposerStaticInitcfb7004a461833b2c554ede0980f1246::getInitializer($loader));

        $loader->register(true);

        return $loader;
    }
}

<?php

// autoload_static.php @generated by Composer

namespace Composer\Autoload;

class ComposerStaticInitcfb7004a461833b2c554ede0980f1246
{
    public static $prefixLengthsPsr4 = array (
        'B' => 
        array (
            'BenWatson\\MyPhpProject\\' => 23,
        ),
    );

    public static $prefixDirsPsr4 = array (
        'BenWatson\\MyPhpProject\\' => 
        array (
            0 => __DIR__ . '/../..' . '/src',
        ),
    );

    public static $classMap = array (
        'Composer\\InstalledVersions' => __DIR__ . '/..' . '/composer/InstalledVersions.php',
    );

    public static function getInitializer(ClassLoader $loader)
    {
        return \Closure::bind(function () use ($loader) {
            $loader->prefixLengthsPsr4 = ComposerStaticInitcfb7004a461833b2c554ede0980f1246::$prefixLengthsPsr4;
            $loader->prefixDirsPsr4 = ComposerStaticInitcfb7004a461833b2c554ede0980f1246::$prefixDirsPsr4;
            $loader->classMap = ComposerStaticInitcfb7004a461833b2c554ede0980f1246::$classMap;

        }, null, ClassLoader::class);
    }
}
